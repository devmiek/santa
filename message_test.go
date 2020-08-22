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

	"github.com/stretchr/testify/assert"
)

func TestStringMessage(t *testing.T) {
	buffer := make([]byte, 0, 256)

	message := StringMessage("Hello Test!")
	buffer = message.FormatStandard(buffer)
	
	assert.Equal(t, "Hello Test!", string(buffer),
		"Unexpected format result")

	buffer = message.FormatJSON(buffer[ : 0])

	assert.Equal(t, `"Hello Test!"`, string(buffer),
		"Unexpected format result")

	assert.Equal(t, "Hello Test!", message.SampleText(),
		"Unexpected sample result")
}

func TestTemplateMessage(t *testing.T) {
	buffer := make([]byte, 0, 256)

	message := TemplateMessage {
		Template: "Hello %s!",
		Args: []interface { } {
			"Test",
		},
	}

	buffer = message.FormatStandard(buffer)

	assert.Equal(t, "Hello Test!", string(buffer),
		"Unexpected format result")

	buffer = message.FormatJSON(buffer[ : 0])

	assert.Equal(t, `"Hello Test!"`, string(buffer),
		"Unexpected format result")

	assert.Equal(t, "Hello %s!", message.SampleText(),
		"Unexpected sample result")
}

func TestStructMessage(t *testing.T) {
	buffer := make([]byte, 0, 256)

	message := StructMessage {
		Text: "Hello Test!",
		Fields: ElementObject {
			String("name", "test"),
			Int("age", 100),
		},
	}

	buffer = message.FormatStandard(buffer)

	assert.Equal(t, `Hello Test! {"name": "test", "age": 100}`,
		string(buffer), "Unexpected format result")

	buffer = message.FormatJSON(buffer[ : 0])

	assert.JSONEq(t, `{
		"textPayload": "Hello Test!",
		"jsonPayload": {
			"name": "test",
			"age": 100
		}
	}`, string(buffer), "Unexpected format result")

	assert.Equal(t, "Hello Test!", message.SampleText(),
		"Unexpected sample result")
}
