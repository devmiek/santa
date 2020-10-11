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

func TestStructMessagePool(t *testing.T) {
	pool := NewStructMessagePool()

	sample := StructMessage {
		Text: "Hello Test!",
		Fields: []Field {
			String("name", "test"),
			Int("age", 100),
		},
	}

	message := pool.New(sample.Text, sample.Fields)

	assert.NotNil(t, message, "Unexpected new error")
	assert.IsType(t, &StructMessage { }, message, "Unexpected new result")
	assert.Equal(t, sample, *message, "Unexpected message value")

	pool.Free(message)
}

func TestTemplateMessagePool(t *testing.T) {
	pool := NewTemplateMessagePool()

	sample := TemplateMessage {
		Template: "Hello %s!",
		Args: []interface { } { "test" },
	}

	message := pool.New(sample.Template, sample.Args)

	assert.NotNil(t, message, "Unexpected new error")
	assert.IsType(t, &TemplateMessage { }, message, "Unexpected new result")
	assert.Equal(t, sample, *message, "Unexpected message value")

	pool.Free(message)
}

func TestEntryPool(t *testing.T) {
	pool := NewEntryPool()

	entry := pool.New()

	assert.NotNil(t, entry, "Unexpected new error")
	assert.IsType(t, &Entry { }, entry, "Unexpected new result")

	pool.Free(entry)
}

func TestExporterBufferPool(t *testing.T) {
	pool := NewExporterBufferPool(2048)

	pointer := pool.New()

	buffer := make([]byte, 0, 1)

	assert.NotNil(t, pointer, "Unexpected new error")
	assert.IsType(t, &buffer, pointer, "Unexpected new result")

	pool.Free(pointer)
}
