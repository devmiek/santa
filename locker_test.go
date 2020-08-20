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
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSpinLock(t *testing.T) {
	locker := NewSpinLock()
	values := make([]uint8, 0)
	waitGroup := &sync.WaitGroup { }

	handler := func(times int) {
		defer waitGroup.Done()

		for count := 0; count < times; count++ {
			locker.Lock()
			values = append(values, 1)
			locker.Unlock()
		}
	}

	const times = 10000

	for count := 0; count < runtime.NumCPU(); count++ {
		waitGroup.Add(1)
		go handler(times)
	}

	go func() {
		locker.LockAndSuspend()

		assert.False(t, locker.TryLock(), "Unexpectedly lock")
		time.Sleep(time.Millisecond * 100)

		locker.UnlockAndResume()
	}()

	waitGroup.Wait()

	assert.Len(t, values, times * runtime.NumCPU(),
		"Unexpected number of elements")
}
