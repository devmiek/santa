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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevelFormatAndParse(t *testing.T) {
	buffer := make([]byte, 0, 32)

	for _, sample := range []struct {
		name string
		level Level
	} {
		{
			name: "debug",
			level: LevelDebug,
		},
		{
			name: "info",
			level: LevelInfo,
		},
		{
			name: "warning",
			level: LevelWarning,
		},
		{
			name: "error",
			level: LevelError,
		},
		{
			name: "fatal",
			level: LevelFatal,
		},
	} {
		assert.Equal(t, sample.name, sample.level.String(),
			"Unexpected level name")

		assert.Equal(t, strings.ToUpper(sample.name), string(
			sample.level.AppendFormat(buffer[ : 0])),
			"Unexpected level format result")
	
		assert.Equal(t, strings.ToUpper(sample.name),
			sample.level.Format(), "Unexpected level format result")

		level, err := ParseLevel(sample.name)
		assert.NoError(t, err, "Unexpected level name parse error")

		assert.Equal(t, sample.level, level,
			"Unexpected level name parse result")
	}
}

func TestLevelEnabled(t *testing.T) {
	for _, sample := range []struct {
		enable Level
		actual Level
		expected bool
	} {
		{
			enable: LevelInfo,
			actual: LevelDebug,
			expected: false,
		},
		{
			enable: LevelDebug,
			actual: LevelDebug,
			expected: true,
		},
		{
			enable: LevelError,
			actual: LevelFatal,
			expected: true,
		},
		{
			enable: LevelFatal,
			actual: LevelWarning,
			expected: false,
		},
	} {
		assert.Equal(t, sample.expected, sample.enable.Enabled(
			sample.actual), "Unexpected result")
	}
}

func TestLevelSpanContains(t *testing.T) {
	for _, sample := range []struct {
		span LevelSpan
		actual Level
		expected bool
	} {
		{
			span: LevelSpan {
				Start: LevelInfo,
				End: LevelWarning,
			},
			actual: LevelFatal,
			expected: false,
		},
		{
			span: LevelSpan {
				Start: LevelError,
				End: LevelFatal,
			},
			actual: LevelError,
			expected: true,
		},
		{
			span: LevelSpan {
				Start: LevelError,
				End: LevelFatal,
			},
			actual: LevelFatal,
			expected: true,
		},
		{
			span: LevelSpan {
				Start: LevelWarning,
				End: LevelError,
			},
			actual: LevelDebug,
			expected: false,
		},
	} {
		assert.Equal(t, sample.expected, sample.span.Contains(
			sample.actual), "Unexpected result")
	}
}
