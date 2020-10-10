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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	entry = &Entry {
		Level: LevelInfo,
		Message: StringMessage("Hello Test!"),
		SourceLocation: EntrySourceLocation {
			File: "main.go",
			Line: 100,
			Parsed: true,
		},
		Name: "test",
		Labels: NewSerializedLabels(
			NewLabel("instanceId", "d325ef24327c"),
		),
	}

	encoderOption = EncoderOption {
		EncodeTime: true,
		EncodeSourceLocation: true,
		EncodeLabels: true,
		EncodeName: true,
		EncodeLevel: true,
	}
)

func init() {
	timestamp, _ := time.Parse(time.RFC3339Nano,
		"2020-08-13T21:56:30.0719939+08:00")

	entry.Time = timestamp
}

func TestStandardEncoderEncode(t *testing.T) {
	buffer := make([]byte, 0, 1024)

	encoder, err := NewStandardEncoder()
	assert.NoError(t, err, "Unexpected standard encoder creation error")
	
	buffer, err = encoder.Encode(buffer, entry)
	assert.NoError(t, err, "Unexpected standard encoder error")

	var expected = fmt.Sprintf("%s %s:%d %s %s [%s] %s\n",
		entry.Time.Format(time.RFC3339Nano),
		entry.SourceLocation.File,
		entry.SourceLocation.Line,
		string(entry.Labels.SerializeStandard(nil)),
		entry.Name,
		entry.Level.Format(),
		entry.Message.(StringMessage),
	)

	assert.Equal(t, expected, string(buffer),
		"Unexpected standard encoder output")
}

func TestJSONEncoderEncode(t *testing.T) {
	buffer := make([]byte, 0, 1024)

	encoder, err := NewJSONEncoder()
	assert.NoError(t, err, "Unexpected JSON encoder creation error")

	buffer, err = encoder.Encode(buffer, entry)
	assert.NoError(t, err, "Unexpected JSON encoder error")

	const expected = `{
		"timestamp": 1597326990071993900,
		"sourceLocation": {
			"file": "main.go",
			"line": 100,
			"function": ""
		},
		"labels": {
			"instanceId": "d325ef24327c"
		},
		"name": "test",
		"level": "INFO",
		"message": "Hello Test!"
	}`

	assert.JSONEq(t, expected, string(buffer),
		"Unexpected JSON encoder output")
}

func TestStandardEncoderOption(t *testing.T) {
	option := NewStandardEncoderOption()

	option.UseTimeLayout(time.RFC3339Nano)
	option.UseEncoderOption(encoderOption)

	assert.Equal(t, time.RFC3339Nano, option.TimeLayout,
		"Unexpected option value")
	
	assert.Equal(t, option.EncoderOption, encoderOption,
		"Unexpected option value")

	_, err := option.Build()
	assert.NoError(t, err, "Unexpected build error")
}

func TestJSONEncoderOption(t *testing.T) {
	option := NewJSONEncoderOption()

	option.UseEncoderOption(encoderOption)
	option.UseTimeLayout(time.RFC3339Nano)

	assert.Equal(t, time.RFC3339Nano, option.TimeLayout,
		"Unexpected option value")
	
	assert.Equal(t, option.EncoderOption, encoderOption,
		"Unexpected option value")

	encoderKeys := NewEncoderKeys()
	encoderKeys.TimeKey = "date"

	option.UseEncoderKeys(encoderKeys)

	assert.Equal(t, encoderKeys, option.EncoderKeys,
		"Unexpected option value")

	_, err := option.Build()
	assert.NoError(t, err, "Unexpected build error")
}
