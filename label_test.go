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

func TestLabelSerialize(t *testing.T) {
	labels := Labels {
		NewLabel("projectId", "santa-project"),
		NewLabel("zoneId", "ap-shanghai-1"),
		NewLabel("instanceId", "d325ef24327c"),
	}

	buffer := make([]byte, 0, 1024)
	buffer = labels.SerializeStandard(buffer)

	assert.JSONEq(t, `{
		"projectId": "santa-project",
		"zoneId": "ap-shanghai-1",
		"instanceId": "d325ef24327c"
	}`, string(buffer), "Unexpected JSON serialization result")
}

func TestSerializedLabels(t *testing.T) {
	labels := NewSerializedLabels(
		NewLabel("projectId", "santa-project"),
		NewLabel("zoneId", "ap-shanghai-1"),
		NewLabel("instanceId", "d325ef24327c"),
	)

	buffer := make([]byte, 0, 1024)
	buffer = labels.SerializeStandard(buffer)

	assert.JSONEq(t, `{
		"projectId": "santa-project",
		"zoneId": "ap-shanghai-1",
		"instanceId": "d325ef24327c"
	}`, string(buffer), "Unexpected JSON serialization result")
}
