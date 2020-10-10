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

// Hook is the public interface of Hook.
//
// Hook is an event callback mechanism. Any Hook type instance that
// implements this interface can be bound to one or more logger instances.
// When one or more log instances bound to the Hook instance trigger an
// event, the corresponding event processing function on the Hook instance
// will be called.
//
// Using the Hook mechanism, developers can intercept and process interesting
// events and log entries.
type Hook interface {
	// Print handles the printed log entries. This function will
	// print the log entry in the bound logger instance currently
	// being called.
	//
	// If the function returns an error, the printing operation for
	// the given log entry will be cancelled.
	//
	// Hook instances can modify log entries during this process.
	Print(entry *Entry) error
}

// SimpleHookHandler is the type of handler function of simple Hook.
type SimpleHookHandler func(entry *Entry) error

// SimpleHook is a structure that contains a Hook processing function.
//
// Simple Hook binds an external handler function to process the log
// entry printing event triggered on the Hook instance, and other events
// will be ignored.
type SimpleHook struct {
	handler SimpleHookHandler
}

// NewSimpleHook creates and returns a simple Hook instance.
func NewSimpleHook(handler SimpleHookHandler) *SimpleHook {
	if handler == nil {
		return nil
	}
	return &SimpleHook {
		handler: handler,
	}
}

// Print handles the printing of log entries on the logger instance,
// but the real processor is the processing function bound to the
// simple Hook instance.
func (h *SimpleHook) Print(entry *Entry) error {
	return h.handler(entry)
}
