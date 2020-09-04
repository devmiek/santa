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

import "sync"

// StructMessagePool is a structure that contains instances of
// cached structured messages.
//
// The structure message pool allows the allocated structure message
// instance to be cached in the pool after use and reused in multiple
// hyper-threading contexts, which will significantly reduce the number
// of heap memory allocations.
type StructMessagePool struct {
	pool *sync.Pool
}

// New gets and returns a reusable message instance from the buffer pool.
// If not, then allocate and return a new message instance.
func (p *StructMessagePool) New(text string, fields []Field) *StructMessage {
	message := p.pool.Get().(*StructMessage)
	message.Text = text
	message.Fields = fields

	return message
}

// Free returns the given message instance to the buffer pool. After the
// refund, the message instance is not allowed to be used again, otherwise
// the behavior is undefined.
func (p *StructMessagePool) Free(message *StructMessage) {
	p.pool.Put(message)
}

// NewStructMessagePool creates and returns a structured message buffer
// pool instance.
func NewStructMessagePool() *StructMessagePool {
	return &StructMessagePool {
		pool: &sync.Pool {
			New: func() interface { } {
				return &StructMessage { }
			},
		},
	}
}

// TemplateMessagePool is a structure that contains instances of
// cached template messages.
//
// The template message pool allows the allocated template message
// instance to be cached in the pool after use and reused in multiple
// hyper-threading contexts, which will significantly reduce the number
// of heap memory allocations.
type TemplateMessagePool struct {
	pool *sync.Pool
}

// New gets and returns a reusable message instance from the buffer pool.
// If not, then allocate and return a new message instance.
func (p *TemplateMessagePool) New(template string, args []interface { }) *TemplateMessage {
	message := p.pool.Get().(*TemplateMessage)
	message.Template = template
	message.Args = args

	return message
}

// Free returns the given message instance to the buffer pool. After the
// refund, the message instance is not allowed to be used again, otherwise
// the behavior is undefined.
func (p *TemplateMessagePool) Free(message *TemplateMessage) {
	p.pool.Put(message)
}

// NewTemplateMessagePool creates and returns a template message buffer
// pool instance.
func NewTemplateMessagePool() *TemplateMessagePool {
	return &TemplateMessagePool {
		pool: &sync.Pool {
			New: func() interface { } {
				return &TemplateMessage { }
			},
		},
	}
}

// EntryPool is a structure that contains instances of cached log entries.
//
// The log entry pool allows the allocated and used log entry instances to
// be cached in the pool for use by other hyper-threading contexts, which
// will significantly reduce the number of heap memory allocations.
//
// Note that any instance of log entry should use this pool allocation.
type EntryPool struct {
	pool *sync.Pool
}

// New gets and returns a reusable log entry instance from the buffer pool.
// If not, then allocate and return a new log entry instance.
//
// Please note that the log entry instance obtained and returned may be dirty,
// and the pool is not responsible for cleaning it.
func (p *EntryPool) New() *Entry {
	return p.pool.Get().(*Entry)
}

// Free returns the given log entry instance to the buffer pool. After the
// refund, the log entry instance is not allowed to be used again, otherwise
// the behavior is undefined.
func (p *EntryPool) Free(entry *Entry) {
	p.pool.Put(entry)
}

// NewEntryPool creates and returns a log entry buffer pool instance.
func NewEntryPool() *EntryPool {
	return &EntryPool {
		pool: &sync.Pool {
			New: func() interface { } {
				return &Entry { }
			},
		},
	}
}

// ExporterBufferPool is a structure that contains instances of cached
// exporter buffer.
//
// The exporter buffer pool allows the allocated and used exporter buffer
// instances to be cached in the pool for use by other hyper-threading
// contexts, which will significantly reduce the number of heap memory
// allocations.
//
// Note that any instance of exporter buffer should use this pool
// allocation.
type ExporterBufferPool struct {
	pool *sync.Pool
}

// New gets and returns a reusable exporter buffer instance from the
// buffer pool. If not, then allocate and return a new exporter buffer
// instance.
//
// Please note that the exporter buffer instance obtained and returned
// may be dirty, and the pool is not responsible for cleaning it.
func (p *ExporterBufferPool) New() *[]byte {
	return p.pool.Get().(*[]byte)
}

// Free returns the given exporter buffer instance to the buffer pool.
// After the refund, the exporter buffer instance is not allowed to be
// used again, otherwise the behavior is undefined.
func (p *ExporterBufferPool) Free(buffer *[]byte) {
	p.pool.Put(buffer)
}

// NewExporterBufferPool creates and returns a log entry buffer pool
// instance.
func NewExporterBufferPool() *ExporterBufferPool {
	return &ExporterBufferPool {
		pool: &sync.Pool {
			New: func() interface { } {
				buffer := make([]byte, 0, 2048)
				return &buffer
			},
		},
	}
}

// DecoratorPool is a structure that contains instances of cached
// decorators.
//
// The decorator pool allows the allocated and used decorator instances
// to be cached in the pool for reuse by other hyper-threading contexts,
// which will significantly reduce the number of heap memory allocations.
type DecoratorPool struct {
	pool *sync.Pool
}

// New gets, initializes and returns a usable decorator instance from the
// buffer pool. If no decorator instance is available, a new one is
// created.
func (p *DecoratorPool) New(logger *Logger) *Decorator {
	decorator := p.pool.Get().(*Decorator)
	decorator.Logger = *logger
	decorator.brush.logger = &decorator.Logger
	return decorator
}

// Free returns the instance of the given decorator to the buffer pool.
// Note that the behavior of using instances of returned decorators is
// undefined.
func (p *DecoratorPool) Free(decorator *Decorator) {
	p.pool.Put(decorator)
}

// NewDecoratorPool creates and returns a decorator pool instance.
func NewDecoratorPool() *DecoratorPool {
	return &DecoratorPool {
		pool: &sync.Pool {
			New: func() interface { } {
				return &Decorator { }
			},
		},
	}
}

// StandardDecoratorPool is a structure that contains instances of
// cached decorators.
//
// The decorator pool allows the allocated and used decorator instances
// to be cached in the pool for reuse by other hyper-threading contexts,
// which will significantly reduce the number of heap memory allocations.
type StandardDecoratorPool struct {
	pool *sync.Pool
}

// New gets, initializes and returns a usable decorator instance from the
// buffer pool. If no decorator instance is available, a new one is
// created.
func (p *StandardDecoratorPool) New(logger *StandardLogger) *StandardDecorator {
	decorator := p.pool.Get().(*StandardDecorator)
	decorator.StandardLogger = *logger
	decorator.brush.logger = &decorator.Logger
	return decorator
}

// Free returns the instance of the given decorator to the buffer pool.
// Note that the behavior of using instances of returned decorators is
// undefined.
func (p *StandardDecoratorPool) Free(decorator *StandardDecorator) {
	p.pool.Put(decorator)
}

// NewStandardDecoratorPool creates and returns a standard decorator
// pool instance.
func NewStandardDecoratorPool() *StandardDecoratorPool {
	return &StandardDecoratorPool {
		pool: &sync.Pool {
			New: func() interface { } {
				return &StandardDecorator { }
			},
		},
	}
}

// TemplateDecoratorPool is a structure that contains instances of
// cached decorators.
//
// The decorator pool allows the allocated and used decorator instances
// to be cached in the pool for reuse by other hyper-threading contexts,
// which will significantly reduce the number of heap memory allocations.
type TemplateDecoratorPool struct {
	pool *sync.Pool
}

// New gets, initializes and returns a usable decorator instance from the
// buffer pool. If no decorator instance is available, a new one is
// created.
func (p *TemplateDecoratorPool) New(logger *TemplateLogger) *TemplateDecorator {
	decorator := p.pool.Get().(*TemplateDecorator)
	decorator.TemplateLogger = *logger
	decorator.brush.logger = &decorator.Logger
	return decorator
}

// Free returns the instance of the given decorator to the buffer pool.
// Note that the behavior of using instances of returned decorators is
// undefined.
func (p *TemplateDecoratorPool) Free(decorator *TemplateDecorator) {
	p.pool.Put(decorator)
}

// NewTemplateDecoratorPool creates and returns a template decorator
// pool instance.
func NewTemplateDecoratorPool() *TemplateDecoratorPool {
	return &TemplateDecoratorPool {
		pool: &sync.Pool {
			New: func() interface { } {
				return &TemplateDecorator { }
			},
		},
	}
}

// StructDecoratorPool is a structure that contains instances of
// cached decorators.
//
// The decorator pool allows the allocated and used decorator instances
// to be cached in the pool for reuse by other hyper-threading contexts,
// which will significantly reduce the number of heap memory allocations.
type StructDecoratorPool struct {
	pool *sync.Pool
}

// New gets, initializes and returns a usable decorator instance from the
// buffer pool. If no decorator instance is available, a new one is
// created.
func (p *StructDecoratorPool) New(logger *StructLogger) *StructDecorator {
	decorator := p.pool.Get().(*StructDecorator)
	decorator.StructLogger = *logger
	decorator.brush.logger = &decorator.Logger
	return decorator
}

// Free returns the instance of the given decorator to the buffer pool.
// Note that the behavior of using instances of returned decorators is
// undefined.
func (p *StructDecoratorPool) Free(decorator *StructDecorator) {
	p.pool.Put(decorator)
}

// NewStructDecoratorPool creates and returns a struct decorator
// pool instance.
func NewStructDecoratorPool() *StructDecoratorPool {
	return &StructDecoratorPool {
		pool: &sync.Pool {
			New: func() interface { } {
				return &StructDecorator { }
			},
		},
	}
}

// pool is a structural variable that contains default instances of
// various pools. These pool instances are automatically created when
// the application is initialized and shared globally.
var pool struct {
	entry *EntryPool
	message struct {
		structure *StructMessagePool
		template *TemplateMessagePool
	}
	buffer struct {
		exporter *ExporterBufferPool
	}
	decorator struct {
		base *DecoratorPool
		standard *StandardDecoratorPool
		template *TemplateDecoratorPool
		structure *StructDecoratorPool
	}
}

// init is used to initialize the variable pool.
func init() {
	pool.entry = NewEntryPool()
	pool.message.template = NewTemplateMessagePool()
	pool.message.structure = NewStructMessagePool()
	pool.buffer.exporter = NewExporterBufferPool()
	pool.decorator.base = NewDecoratorPool()
	pool.decorator.standard = NewStandardDecoratorPool()
	pool.decorator.template = NewTemplateDecoratorPool()
	pool.decorator.structure = NewStructDecoratorPool()
}
