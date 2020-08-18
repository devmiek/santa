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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStandardSyncerOption(t *testing.T) {
	option := NewStandardSyncerOption()

	option.UseWriter(os.Stderr)
	option.UseCacheCapacity(0)

	assert.Equal(t, os.Stderr, option.Writer, "Unexpected option value")
	assert.Equal(t, 0, option.CacheCapacity, "Unexpected option value")

	syncer, err := option.Build()

	assert.NoError(t, err, "Unexpected build error")
	assert.NotNil(t, syncer, "Unexpected build result")

	assert.Equal(t, option.Writer, syncer.writer,
		"Unexpected instance error")
	assert.Equal(t, option.CacheCapacity, syncer.capacity,
		"Unexpected instance error")

	assert.NoError(t, syncer.Close(), "Unexpected close error")
}

func TestFileSyncerOption(t *testing.T) {
	option := NewFileSyncerOption()

	option.UseName(os.DevNull)
	option.UseCacheCapacity(256)  // Invalid value

	assert.Equal(t, os.DevNull, option.FileName, "Unexpected option value")
	assert.Equal(t, 256, option.CacheCapacity, "Unexpected option value")

	syncer, err := option.Build()
	
	assert.NoError(t, err, "Unexpected build error")
	assert.NotNil(t, syncer, "Unexpected build result")

	assert.NotNil(t, syncer.writer, "Unexpected instance error")
	assert.Equal(t, 1024, syncer.capacity, "Unexpected instance error")
	
	assert.NoError(t, syncer.Close(), "Unexpected close error")
}

func TestStandardSyncerWrite(t *testing.T) {
	syncer, err := NewStandardSyncer()
	assert.NoError(t, err, "Unexpected create error")

	for count := 0; count < 100000; count++ {
		_, err = syncer.Write([]byte("Hello Test!"))
		assert.NoError(t, err, "Unexpected write error")
	}

	assert.NoError(t, syncer.Sync(), "Unexpected sync error")
	assert.NoError(t, syncer.Close(), "Unexpected close error")

	syncer, err = NewStandardSyncerOption().UseCacheCapacity(0).Build()
	assert.NoError(t, err, "Unexpected create error")

	for count := 0; count < 100; count++ {
		_, err = syncer.Write([]byte("Hello Test!"))
		assert.NoError(t, err, "Unexpected write error")
	}

	assert.NoError(t, syncer.Sync(), "Unexpected sync error")
	assert.NoError(t, syncer.Close(), "Unexpected close error")
}

func TestFileSyncerWrite(t *testing.T) {
	syncer, err := NewFileSyncer()
	assert.NoError(t, err, "Unexpected create error")

	_, err = syncer.Write([]byte("Hello Test!"))
	assert.NoError(t, err, "Unexpected write error")
	assert.NoError(t, syncer.Close(), "Unexpected close error")
}
