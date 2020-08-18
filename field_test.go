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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestElementFormatJSON(t *testing.T) {
	timestamp, _ := time.Parse(time.RFC3339Nano,
		"2020-08-13T21:56:30.0719939+08:00")

	buffer := make([]byte, 0, 256)

	for _, sample := range []struct {
		name string
		field Field
		expected string
	} {
		{
			name: "int",
			field: Int("int", 10),
			expected: "10",
		},
		{
			name: "uint",
			field: Uint("uint", 20),
			expected: "20",
		},
		{
			name: "float32",
			field: Float32("float32", 3.14),
			expected: "3.14",
		},
		{
			name: "float64",
			field: Float64("float64", 3.1415),
			expected: "3.1415",
		},
		{
			name: "boolean",
			field: Boolean("boolean", true),
			expected: "true",
		},
		{
			name: "string",
			field: String("string", "Hello"),
			expected: "\"Hello\"",
		},
		{
			name: "bytes",
			field: Bytes("bytes", []byte("Hello")),
			expected: "\"Hello\"",
		},
		{
			name: "time",
			field: Time("time", timestamp),
			expected: "1597326990071993900",
		},
		{
			name: "error",
			field: Error("error", errors.New("Error")),
			expected: "\"Error\"",
		},
		{
			name: "value",
			field: Value("value", 50),
			expected: "50",
		},
	} {
		assert.Equal(t, sample.name, sample.field.Name,
			"Unexpected field name")

		assert.Equal(t, sample.expected, string(
			sample.field.FormatJSON(buffer[ : 0])),
			"Unexpected JSON formatted append result",
		)
	}

	fields := []Field {
		String("name", "test"),
		Int("age", 100),
	}
	
	for _, sample := range []struct {
		name string
		field Field
		expected string
	} {
		{
			name: "object",
			field: Object("object", fields...),
			expected: `{
				"name": "test",
				"age": 100
			}`,
		},
		{
			name: "objects",
			field: Objects("objects",
				ElementObject(fields),
				ElementObject(fields),
			),
			expected: `[
				{
					"name": "test",
					"age": 100
				},
				{
					"name": "test",
					"age": 100
				}
			]`,
		},
		{
			name: "ints",
			field: Ints("ints", []int64 { 10, 20, 30 }),
			expected: `[10, 20, 30]`,
		},
		{
			name: "uints",
			field: Uints("uints", []uint64 { 40, 50, 60 }),
			expected: `[40, 50, 60]`,
		},
		{
			name: "float32s",
			field: Float32s("float32s", []float32 { 1.1, 1.2,
				1.3 }),
			expected: `[1.1, 1.2, 1.3]`,
		},
		{
			name: "float64s",
			field: Float64s("float64s", []float64 { 1.4, 1.5, 1.6 }),
			expected: `[1.4, 1.5, 1.6]`,
		},
		{
			name: "booleans",
			field: Booleans("booleans", []bool { true, false }),
			expected: `[true, false]`,
		},
		{
			name: "strings",
			field: Strings("strings", []string { "value1", "value2" }),
			expected: `["value1", "value2"]`,
		},
		{
			name: "times",
			field: Times("times", []time.Time { timestamp,
				timestamp }),
			expected: `[1597326990071993900, 1597326990071993900]`,
		},
	} {
		assert.Equal(t, sample.name, sample.field.Name,
			"Unexpected field name")

		assert.JSONEq(t, sample.expected,
			string(sample.field.FormatJSON(buffer[ : 0])),
			"Unexpected JSON formatted append result",
		)
	}
}
