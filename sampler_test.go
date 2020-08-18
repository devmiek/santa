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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTextSamplerOption(t *testing.T) {
	option := NewTextSamplerOption()

	span := LevelSpan {
		Start: LevelInfo,
		End: LevelError,
	}

	option.UseSpan(span.Start, span.End)
	option.UseTick(time.Second * 2)
	option.UseFirst(50, 100)
	option.UseCounters(2048)

	assert.Equal(t, span, option.Span, "Unexpected option value")
	assert.Equal(t, time.Second * 2, option.Tick, "Unexpected option value")
	assert.Equal(t, uint64(50), option.First, "Unexpected option value")
	assert.Equal(t, uint64(100), option.Thereafter, "Unexpected option value")
	assert.Equal(t, uint64(2048), option.Counters, "Unexpected option value")

	sampler, err := option.Build()
	assert.NoError(t, err, "Unexpected build error")

	assert.Equal(t, option.Span, sampler.span, "Unexpected instance error")
	assert.Equal(t, int64(option.Tick), sampler.tick,
		"Unexpected instance error")
	assert.Equal(t, option.First, sampler.first, "Unexpected instance error")
	assert.Equal(t, option.Thereafter, sampler.thereafter,
		"Unexpected instance error")
	assert.Equal(t, int(option.Counters), len(sampler.counters),
		"Unexpected instance error")
}

type testMessage struct { }

func TestTextSamplerSample(t *testing.T) {
	sampler, err := NewTextSampler()
	assert.NoError(t, err, "Unexpected create error")

	for _, entry := range []Entry {
		{
			Time: time.Now(),
			Level: LevelInfo,
			Message: StringMessage("Hello Test!"),
		},
		{
			Time: time.Now(),
			Level: LevelFatal,
			Message: StringMessage("Hello Test!"),
		},
		{
			Time: time.Now(),
			Level: LevelInfo,
			Message: testMessage { },
		},
	} {
		assert.True(t, sampler.Sample(&entry), "Unexpected sampling result")
	}

	entry := Entry {
		Time: time.Now(),
		Level: LevelInfo,
		Message: StringMessage("Hello Test!"),
	}

	for count := 0; count < 1000; count++ {
		if !sampler.Sample(&entry) {
			break
		}

		if count == 1000 {
			assert.Fail(t, "Sampler not working")
		}
	}
}
