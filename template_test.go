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

func TestTemplateLoggerOption(t *testing.T) {
	option := NewTemplateOption()

	encodingOption := NewEncodingOption()
	samplingOption := NewSamplingOption()
	outputtingOption := NewOutputtingOption()
	flushingOption := NewFlushingOption()

	option.UseEncoding(encodingOption)
	option.UseSampling(samplingOption)
	option.UseOutputting(outputtingOption)
	option.UseErrorOutputting(outputtingOption)
	option.UseFlushing(flushingOption)

	assert.Equal(t, *encodingOption, option.Encoding,
		"Unexpected option value")
	assert.Equal(t, *samplingOption, option.Sampling,
		"Unexpected option value")
	assert.Equal(t, *outputtingOption, option.Outputting,
		"Unexpected option value")
	assert.Equal(t, *outputtingOption, option.ErrorOutputting,
		"Unexpected option value")
	assert.Equal(t, *flushingOption, option.Flushing,
		"Unexpected option value")

	option.UseHook(NewSimpleHook(func(entry *Entry) error {
		return nil
	}))

	option.UseLabel(NewLabel("instanceId", "d325ef24327c"))
	option.UseLevel(LevelInfo)
	option.UseName("test")
	
	assert.Equal(t, "d325ef24327c", option.Labels[0].Value)
	assert.Equal(t, LevelInfo, option.Level, "Unexpected option value")
	assert.Equal(t, "test", option.Name, "Unexpected option value")

	logger, err := option.Build()
	assert.NoError(t, err, "Unexpected build error")

	assert.NotNil(t, logger.sampler, "Unexpected instance error")
	assert.Len(t, logger.exporters, 2, "Unexpected instance error")
	assert.NotNil(t, logger.exporters[0], "Unexpected instance error")
	assert.NotNil(t, logger.exporters[1], "Unexpected instance error")

	assert.Equal(t, option.Level, logger.level, "Unexpected instance error")
	assert.Equal(t, option.Name, logger.name, "Unexpected instance error")

	option.DisableCache()
	option.DisableFlushing()
	option.DisableSampling()

	assert.Equal(t, SamplingOption { }, option.Sampling,
		"Unexpected option value")
	assert.Equal(t, time.Duration(0), option.Flushing.Interval,
		"Unexpected option value")
	assert.True(t, option.Outputting.DisableCache,
		"Unexpected option value")
	assert.True(t, option.ErrorOutputting.DisableCache,
		"Unexpected option value")

	logger, err = option.Build()
	assert.NoError(t, err, "Unexpected build error")
	assert.NotNil(t, logger, "Unexpected build result")
	assert.NoError(t, logger.Close(), "Unexpected close error")
}

func TestTemplateLoggerBenchmark(t *testing.T) {
	logger, err := NewTemplateBenchmark(true, EncoderJSON)
	assert.NoError(t, err, "Unexpected create error")
	assert.NotNil(t, logger, "Unexpected create error")

	logger, err = NewTemplateBenchmark(false, EncoderJSON)
	assert.NoError(t, err, "Unexpected create error")
	assert.NotNil(t, logger, "Unexpected create error")

	logger, err = NewTemplateBenchmark(true, EncoderStandard)
	assert.NoError(t, err, "Unexpected create error")
	assert.NotNil(t, logger, "Unexpected create error")

	logger, err = NewTemplateBenchmark(false, EncoderStandard)
	assert.NoError(t, err, "Unexpected create error")
	assert.NotNil(t, logger, "Unexpected create error")

	assert.NoError(t, logger.Close(), "Unexpected close error")
}

func TestTemplateLoggerPrint(t *testing.T) {
	logger, err := NewTemplate()
	assert.NoError(t, err, "Unexpected create error")
	assert.NoError(t, logger.Sync(), "Unexpected sync error")
	assert.NoError(t, logger.Close(), "Unexpected close error")
	
	logger, err = NewTemplateBenchmark(false, EncoderJSON)
	assert.NoError(t, err, "Unexpected create error")

	err = logger.Debugf("Hello Test! %s %d", "test", 100)
	assert.NoError(t, err, "Unexpected print error")

	err = logger.Infof("Hello Test! %s %d", "test", 100)
	assert.NoError(t, err, "Unexpected print error")

	err = logger.Warningf("Hello Test! %s %d", "test", 100)
	assert.NoError(t, err, "Unexpected print error")

	err = logger.Errorf("Hello Test! %s %d", "test", 100)
	assert.NoError(t, err, "Unexpected print error")

	err = logger.Fatalf("Hello Test! %s %d", "test", 100)
	assert.NoError(t, err, "Unexpected print error")

	err = logger.Printf(LevelError, "Hello Test! %s %d", "test", 100)
	assert.NoError(t, err, "Unexpected print error")

	assert.NoError(t, logger.Close(), "Unexpected close error")
}
