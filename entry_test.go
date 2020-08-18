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

func TestEntrySourceLocation(t *testing.T) {
	buffer := make([]byte, 0, 256)

	sourceLocation := EntrySourceLocation {
		File: "main.go",
		Line: 100,
		Parsed: true,
	}

	buffer = sourceLocation.AppendString(buffer)

	assert.Equal(t, "main.go:100", string(buffer),
		"Unexpected append result")
	
	buffer = sourceLocation.AppendJSON(buffer[ : 0])

	const expected = `{
        "file": "main.go",
        "line": 100,
        "function": ""
	}`

	assert.JSONEq(t, expected, string(buffer),
		"Unexpected append result")
}
