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
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoggerOption(t *testing.T) {
	exporter, _ := NewStandardExporter()
	sampler, _ := NewTextSampler()

	option := NewOption()

	option.Labels = append(option.Labels, NewLabel("instanceId", "d325ef24327c"))
	option.Exporters = append(option.Exporters, exporter)
	option.Sampler = sampler
	option.Level = LevelInfo
	option.Name = "test"

	logger, err := option.Build()
	assert.NoError(t, err, "Unexpected build error")

	assert.Equal(t, 1, logger.labels.Count(), "Unexpected instance error")
	assert.Len(t, logger.exporters, 1, "Unexpected instance error")
	assert.Equal(t, exporter, logger.exporters[0], "Unexpected instance error")
	assert.Equal(t, option.Sampler, logger.sampler, "Unexpected instance error")
	assert.Equal(t, option.Level, logger.level, "Unexpected instance error")
	assert.Equal(t, option.Name, logger.name, "Unexpected instance error")
}

type testExporter struct {
	entry *Entry
}

func (e *testExporter) Export(entry *Entry) error {
	e.entry = entry
	return nil
}

func (e *testExporter) Sync() error {
	return nil
}

func (e *testExporter) Close() error {
	return nil
}

func TestLoggerPrint(t *testing.T) {
	option := NewOption()
	option.Exporters = append(option.Exporters, &testExporter { })

	logger, err := option.Build()
	assert.NoError(t, err, "Unexpected create error")

	err = logger.Print(LevelInfo, StringMessage("Hello Test!"))
	assert.NoError(t, err, "Unexpected print error")

	assert.Equal(t, LevelInfo, option.Exporters[0].(*testExporter).
		entry.Level, "Unexpected log entry")

	assert.Equal(t, StringMessage("Hello Test!"), option.Exporters[0].
		(*testExporter).entry.Message.(StringMessage),
		"Unexpected log entry")
}

func TestEncodingOption(t *testing.T) {
	option := NewEncodingOption()
	option.UseStandard()

	assert.Equal(t, option.Kind, EncoderStandard, "Unexpected option value")
	assert.IsType(t, &StandardEncoderOption { }, option.Option,
		"Unexpected option value")

	encoder, err := option.Build()
	assert.NoError(t, err, "Unexpected build error")

	assert.IsType(t, &StandardEncoder { }, encoder,
		"Unexpected instance error")

	option.UseJSON()

	assert.Equal(t, option.Kind, EncoderJSON, "Unexpected option value")
	assert.IsType(t, &JSONEncoderOption { }, option.Option,
		"Unexpected option value")

	encoder, err = option.Build()
	assert.NoError(t, err, "Unexpected build error")

	assert.IsType(t, &JSONEncoder { }, encoder,
		"Unexpected instance error")

	standardEncoderOption := NewStandardEncoderOption()
	option.UseStandardOption(standardEncoderOption)

	assert.Equal(t, option.Kind, EncoderStandard, "Unexpected option value")
	assert.Equal(t, standardEncoderOption, option.Option,
		"Unexpected option value")
	
	encoder, err = option.Build()
	assert.NoError(t, err, "Unexpected build error")

	assert.IsType(t, &StandardEncoder { }, encoder,
		"Unexpected instance error")

	jsonEncoderOption := NewJSONEncoderOption()
	option.UseJSONOption(jsonEncoderOption)

	assert.Equal(t, option.Kind, EncoderJSON, "Unexpected option value")
	assert.Equal(t, jsonEncoderOption, option.Option,
		"Unexpected option value")

	encoder, err = option.Build()
	assert.NoError(t, err, "Unexpected build error")

	assert.IsType(t, &JSONEncoder { }, encoder,
		"Unexpected instance error")
}

func TestSamplingOption(t *testing.T) {
	option := NewSamplingOption()
	option.UseText()

	textSamplerOption := NewTextSamplerOption()

	assert.Equal(t, SamplerText, option.Kind, "Unexpected option value")
	assert.IsType(t, textSamplerOption, option.Option,
		"Unexpected option value")

	option.UseTextOption(textSamplerOption)

	assert.Equal(t, SamplerText, option.Kind, "Unexpected option value")
	assert.Equal(t, textSamplerOption, option.Option,
		"Unexpected option value")

	sampler, err := option.Build()
	assert.NoError(t, err, "Unexpected build error")

	assert.IsType(t, &TextSampler { }, sampler,
		"Unexpected instance error")
}

func TestOutputtingOption(t *testing.T) {
	option := NewOutputtingOption()

	option.UseStandard(ioutil.Discard)

	assert.Equal(t, SyncerStandard, option.Kind, "Unexpected option value")
	assert.IsType(t, &StandardSyncerOption { }, option.Option,
		"Unexpected option value")
	assert.Equal(t, ioutil.Discard, option.Option.(*StandardSyncerOption).
		Writer, "Unexpected option value")
	
	syncer, err := option.Build()
	assert.NoError(t, err, "Unexpected build error")

	assert.IsType(t, &StandardSyncer { }, syncer,
		"Unexpected instance error")

	assert.NoError(t, syncer.Close(), "Unexpected close error")

	option.UseFile(os.DevNull)

	assert.Equal(t, SyncerFile, option.Kind, "Unexpected option value")
	assert.IsType(t, &FileSyncerOption { }, option.Option,
		"Unexpected option value")
	assert.Equal(t, os.DevNull, option.Option.(*FileSyncerOption).
		FileName, "Unexpected option value")
	
	syncer, err = option.Build()
	assert.NoError(t, err, "Unexpected build error")

	assert.IsType(t, &FileSyncer { }, syncer,
		"Unexpected instance error")

	assert.NoError(t, syncer.Close(), "Unexpected close error")

	option.UseDiscard()

	assert.Equal(t, SyncerDiscard, option.Kind, "Unexpected option value")
	
	syncer, err = option.Build()
	assert.NoError(t, err, "Unexpected build error")

	assert.IsType(t, &DiscardSyncer { }, syncer,
		"Unexpected instance error")

	assert.NoError(t, syncer.Close(), "Unexpected close error")
}

func TestFlushingOption(t *testing.T) {
	option := NewFlushingOption()
	option.UseInterval(time.Minute)

	assert.Equal(t, time.Minute, option.Interval, "Unexpected option value")
}

func TestStandardLoggerOption(t *testing.T) {
	option := NewStandardOption()

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

	option.UseHooks(NewSimpleHook(func(entry *Entry) error {
		return nil
	}))
	
	option.UseLevel(LevelInfo)
	option.UseName("test")
	
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

func TestStandardLoggerBenchmark(t *testing.T) {
	logger, err := NewStandardBenchmark(true, EncoderJSON)
	assert.NoError(t, err, "Unexpected create error")
	assert.NotNil(t, logger, "Unexpected create error")

	logger, err = NewStandardBenchmark(false, EncoderJSON)
	assert.NoError(t, err, "Unexpected create error")
	assert.NotNil(t, logger, "Unexpected create error")

	logger, err = NewStandardBenchmark(true, EncoderStandard)
	assert.NoError(t, err, "Unexpected create error")
	assert.NotNil(t, logger, "Unexpected create error")

	logger, err = NewStandardBenchmark(false, EncoderStandard)
	assert.NoError(t, err, "Unexpected create error")
	assert.NotNil(t, logger, "Unexpected create error")

	assert.NoError(t, logger.Close(), "Unexpected close error")
}

func TestStandardLoggerPrint(t *testing.T) {
	logger, err := NewStandard()
	assert.NoError(t, err, "Unexpected create error")
	assert.NoError(t, logger.Sync(), "Unexpected sync error")
	assert.NoError(t, logger.Close(), "Unexpected close error")
	
	logger, err = NewStandardBenchmark(false, EncoderJSON)
	assert.NoError(t, err, "Unexpected create error")

	err = logger.Debug(StringMessage("Hello Test!"))
	assert.NoError(t, err, "Unexpected print error")

	err = logger.Info(StringMessage("Hello Test!"))
	assert.NoError(t, err, "Unexpected print error")

	err = logger.Warning(StringMessage("Hello Test!"))
	assert.NoError(t, err, "Unexpected print error")

	err = logger.Error(StringMessage("Hello Test!"))
	assert.NoError(t, err, "Unexpected print error")

	err = logger.Fatal(StringMessage("Hello Test!"))
	assert.NoError(t, err, "Unexpected print error")

	assert.NoError(t, logger.Close(), "Unexpected close error")
}
