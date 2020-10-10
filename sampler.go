// MIT License
//
// Copyright (c) 2020 Nobody Night
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package santa

import (
	"sync/atomic"
	"time"
)

// Sampler is the public interface of the sampler.
//
// The sampler is a log entry sampler, which usually collects part of all
// printed log entries according to certain rules, and other uncollected
// log entries will be discarded.
//
// With the help of the sampler, the utilization of CPU resources and I/O
// resources can be reduced when a large number of log entries are printed
// out to reduce the impact on system performance. Conversely, if the sampler
// is not configured properly, important log entries may be missed.
//
// All valid samplers only need to implement this interface.
type Sampler interface {
	// Sample checks whether a given log entry needs to be sampled. It returns
	// true if needed, otherwise it returns false.
	Sample(entry *Entry) bool
}

type textSamplerCounter struct {
	// count represents the value of the counter.
	count uint64

	// after represents the time when the counter will be reset next.
	after int64
}

// TextSampler is the structure of the text sampler instance.
//
// The text sampler determines whether one or more log entries should be
// discarded by tracking the output times and cycles of each log entry
// message that can be parsed into a text string.
//
// The text sampler checks if the output log entry message is the same
// as a log entry message that has been previously output. If they are
// identical, the sampling policy is to allow the same log entry to be
// output N times in a sampling period, and then output it once after
// every N interval. Identical log entries that do not match the sampling
// policy will be discarded.
//
// Note that the text sampler cannot guarantee the accuracy of detection
// of the same log entry message, which means that the false alarm rate
// increases as more samples of different log entry messages are tracked.
// A way to reduce the false alarm rate is to increase the number of
// counters used for tracking, which can be achieved by adjusting the value
// of the TextSamplerOption.Counters option. More counters means more memory
// resources will be used.
type TextSampler struct {
	span LevelSpan
	tick int64
	first uint64
	thereafter uint64
	counters []textSamplerCounter
}

// TextSampleParser is the public interface of the text sample parser.
//
// The text sample parser is used to parse log entry messages into text. Any log
// entry messages that support text samplers should implement this interface,
// otherwise the text sampler does not know how to parse log entry messages.
type TextSampleParser interface {
	// SampleText returns the text sample string of the log entry message.
	SampleText() string
}

// hash64 uses the FNV64-A algorithm to calculate and returns the Hash value
// of the given text.
func (*TextSampler) hash64(text string) uint64 {
	result := uint64(14695981039346656037)
	for index := 0; index < len(text); index++ {
		result ^= uint64(text[index])
		result *= 1099511628211
	}
	return result
}

// Sample checks whether a given log entry needs to be sampled. It returns
// true if needed, otherwise it returns false.
func (s *TextSampler) Sample(entry *Entry) bool {
	if !s.span.Contains(entry.Level) {
		return true
	}
	parser, ok := entry.Message.(TextSampleParser)
	if !ok {
		return true
	}

	index := s.hash64(parser.SampleText()) % uint64(len(s.counters))
	count := atomic.LoadUint64(&s.counters[index].count)
	clock := entry.Time.UnixNano()
	after := atomic.LoadInt64(&s.counters[index].after)
	
	// If it has been more than or equal to one sampling period since the
	// last time the counter was reset, the counter is reset.
	if after <= clock {
		// Update the next reset time to the counter. If the update fails,
		// it is considered that another hyperthread is competing.
		atomic.CompareAndSwapInt64(&s.counters[index].after, after,
			clock + s.tick)
		
		// If the instant counter count is greater than 0, it is reset to 1.
		if count > 0 {
			// Using subtraction to reset the counter value can avoid
			// incorrect value overwriting when multiple hyperthreads
			// compete.
			atomic.AddUint64(&s.counters[index].count, -count + 1)
		}

		return true
	}

	count = atomic.AddUint64(&s.counters[index].count, 1)

	// If the same log entry has been repeatedly printed <s.first> times in
	// a sampling period, and the condition of printing once after the interval
	// <s.thereafter> times is not met, it will be discarded.
	if count > s.first && (count - s.first) % s.thereafter != 0 {
		return false
	}

	return true
}

// TextSamplerOption is a structure containing text sampler options.
type TextSamplerOption struct {
	// Span represents the log level span for which sampling strategy
	// needs to be applied. If the level of the log entry is not included
	// in the span, the output is sampled.
	//
	// If this option is not set, the default is INFO to WARNING.
	Span LevelSpan

	// Tick represents the sampling cycle time, and the sampling counter
	// is reset every other cycle.
	//
	// If this option is not set, the default is 1 second.
	Tick time.Duration

	// First represents how many times the same log entry message should
	// be allowed to be output repeatedly before discarding it.
	//
	// If this option is not provided, the default is 100 times.
	First uint64

	// Thereafter represents how many times the same log entry message is
	// discarded, it should be allowed to be output once.
	//
	// If this option is not provided, the default is 100 times.
	Thereafter uint64

	// Counters represents the number of counters used to track the same
	// log entry messages. More counters means that checking the same log
	// entry messages is more accurate, but will also consume more memory
	// resources.
	//
	// If this option is not provided, the default is 1024 times.
	Counters uint64
}

// Build builds and returns a text sampler instance using the option value.
//
// Please note that this function does not check the validity of the option
// value, please use the NewTextSamplerOption function to create an option
// instance.
func (o *TextSamplerOption) Build() (*TextSampler, error) {
	return &TextSampler {
		span: o.Span,
		tick: int64(o.Tick),
		first: o.First,
		thereafter: o.Thereafter,
		counters: make([]textSamplerCounter, o.Counters),
	}, nil
}

// UseSpan sets the Span option using the given log level span.
func (o *TextSamplerOption) UseSpan(start, end Level) *TextSamplerOption {
	o.Span = LevelSpan {
		Start: start,
		End: end,
	}
	return o
}

// UseTick sets the Tick option using the given sampling period value.
func (o *TextSamplerOption) UseTick(tick time.Duration) *TextSamplerOption {
	o.Tick = tick
	return o
}

// UseFirst sets the options First and Then using the given value.
func (o *TextSamplerOption) UseFirst(first, thereafter uint64) *TextSamplerOption {
	o.First = first
	o.Thereafter = thereafter
	return o
}

// UseCounters sets the Counters option using the given number of
// sampling counters.
func (o *TextSamplerOption) UseCounters(counters uint64) *TextSamplerOption {
	o.Counters = counters
	return o
}

// NewTextSamplerOption creates and returns a text sampler option instance
// with default option values.
func NewTextSamplerOption() *TextSamplerOption {
	return &TextSamplerOption {
		Span: LevelSpan {
			Start: LevelInfo,
			End: LevelWarning,
		},
		Tick: time.Second,
		First: 100,
		Thereafter: 100,
		Counters: 1024,
	}
}

// NewTextSampler creates and returns a text sampler instance using default
// option values.
func NewTextSampler() (*TextSampler, error) {
	return NewTextSamplerOption().Build()
}
