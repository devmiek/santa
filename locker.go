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
	"sync/atomic"
)

// SpinLock is the structure of the spin lock instance.
//
// Spin locks implement critical sections in user space. Unlike mutex locks,
// spin locks acquire lock ownership by constantly trying to preempt. The
// preemption is carried out in user space. If the preemption fails once,
// the spinlock will actively abandon the running of the current coroutine
// in order to schedule the processor to run other coroutine logic, and then
// try to preempt again.
//
// It is worth noting that spin locks are not suitable for long-term
// possession of the critical section of lock ownership. This will cause
// other coroutine logic trying to preempt the lock to continuously try to
// obtain the lock ownership, which will occupy more CPU time.
//
// If one or more critical regions may occupy one or more spin locks from
// time to time for a long time, please obtain lock ownership by suspending
// the preemption of the spin locks, and resume the suspension when leaving
// the critical region.
//
// The API provided by the spinlock is thread-safe.
type SpinLock struct {
	status int32
	condition *sync.Cond
}

// TryLock attempts to obtain ownership of the lock. It returns true if the
// lock ownership is successfully obtained, otherwise it returns false.
func (l *SpinLock) TryLock() bool {
	return atomic.CompareAndSwapInt32(&l.status, 0, 1)
}

// Lock acquires the ownership of the lock and returns after successfully
// acquiring it.
func (l *SpinLock) Lock() {
	for count := 1; ; count++ {
		if atomic.CompareAndSwapInt32(&l.status, 0, 1) {
			return
		}

		if (count % 100) != 0 {
			runtime.Gosched()
			continue
		}

		l.condition.L.Lock()

		if atomic.LoadInt32(&l.status) == 2 {
			l.condition.Wait()
		}

		l.condition.L.Unlock()
	}
}

// Unlock releases the ownership of the lock.
func (l *SpinLock) Unlock() {
	atomic.StoreInt32(&l.status, 0)
}

// LockAndSuspend obtains the ownership of the lock, then suspends the
// preemption of the ownership of the lock, and returns after success.
func (l *SpinLock) LockAndSuspend() {
	for {
		if atomic.CompareAndSwapInt32(&l.status, 0, 2) {
			break
		}

		runtime.Gosched()
	}
}

// UnlockAndResume releases the lock ownership, and then resumes the
// preemption of the lock ownership.
func (l *SpinLock) UnlockAndResume() {
	l.condition.L.Lock()

	if atomic.CompareAndSwapInt32(&l.status, 2, 0) {
		l.condition.Broadcast()
	}

	l.condition.L.Unlock()
}

// NewSpinLock creates and returns a spinlock instance. For details,
// see the comment section of the SpinLock structure.
func NewSpinLock() *SpinLock {
	return &SpinLock {
		status: 0,
		condition: sync.NewCond(&sync.Mutex { }),
	}
}
