// pool.go is the golang-1.14.6 source file

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

// pool is a structural variable that contains default instances of
// various pools. These pool instances are automatically created when
// the application is initialized and shared globally.
var pool struct {
	entry *EntryPool
	message struct {
		structure *StructMessagePool
		template *TemplateMessagePool
	}
}

// init is used to initialize the variable pool.
func init() {
	pool.entry = NewEntryPool()
	pool.message.template = NewTemplateMessagePool()
	pool.message.structure = NewStructMessagePool()
}
