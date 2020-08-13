// syncer.go is the golang-1.14.6 source file.

package santa

import (
	"io"
	"io/ioutil"
	"os"
)

// Syncer is the public interface of the synchronizer.
//
// The synchronizer writes the encoded log entry data to a specific log
// storage device. Common storage devices include but are not limited to
// standard output, local hard drives, and network locations.
//
// One of the most common scenarios is to use a file synchronizer to write
// encoded log entry data to a local hard disk, where the local hard disk
// is a persistent log storage device.
//
// In general, all types of synchronizers enable internal cache and file
// system cache by default to significantly improve the I/O performance of
// log entry data. However, the side effect of enabling internal caching
// and file system caching is that the time when some log entry data is
// written to the persistent storage device will be delayed, because the
// cache is not flushed.
//
// By default, the internal cache will be flushed automatically after the
// cache capacity is saturated, but the automatic flushing timing of the
// file system cache is determined by the operating system. However, the
// internal cache and file system cache can be flushed manually by calling
// the Sync function.
//
// If you need log entry data to be written to the file system in real
// time, disable the internal cache.
//
// Please note that the synchronizer will automatically flush the internal
// cache and file system cache once when it is closed to ensure that the
// log entry data will not be lost. However, regardless of whether the
// internal cache is disabled, it must be manually closed after the
// synchronizer is no longer in use, otherwise it may cause some log entry
// data loss and file handle leakage.
type Syncer interface {
	// Write writes the data of a given buffer slice to a specific storage
	// device. If the internal cache is enabled, the internal cache is
	// written first. If the capacity of the internal cache is saturated,
	// it is automatically flushed once.
	//
	// Finally, it returns the number of bytes actually written and any
	// errors encountered.
	Write(buffer []byte) (int, error)

	// Sync writes the internally cached data to a specific storage device.
	// If the specific storage device is based on the file system, write the
	// data cached by the file system to the persistent storage device.
	//
	// Finally, any errors encountered are returned.
	Sync() error

	// Close automatically flushes the internal cache once, and then releases
	// any kernel objects that have been opened (including but not limited to:
	// file handles, etc.).
	//
	// Finally, any errors encountered are returned.
	Close() error
}

// SyncerOption is a structure containing basic synchronizer options.
//
// The synchronizer options include basic synchronizer options. Normally,
// all synchronizer option types include this structure.
type SyncerOption struct {
	// CacheCapacity represents the number of bytes of internal cache
	// capacity. If the internal cache is disabled, the cache capacity
	// must be set to 0, otherwise the cache capacity must be greater
	// than or equal to 1,024 bytes.
	//
	// If not provided, the default is 32 KB.
	CacheCapacity int
}

// NewSyncerOption returns the value of a synchronizer option with the
// default optional value.
func NewSyncerOption() SyncerOption {
	return SyncerOption {
		CacheCapacity: 1024 * 32,
	}
}

// StandardSyncer is the structure of a standard synchronizer instance.
//
// The standard synchronizer uses an instance that implements the io.Writer
// interface as a specific storage device, but it does not check the actual
// data type of the instance when it is closed, and it does not close the
// instance (if supported).
//
// Please note that standard synchronizers are not thread-safe.
type StandardSyncer struct {
	writer io.Writer
	buffer []byte
	capacity int
}

// Write writes the data of a given buffer slice to a specific storage
// device. If the internal cache is enabled, the internal cache is
// written first. If the capacity of the internal cache is saturated,
// it is automatically flushed once.
//
// Finally, it returns the number of bytes actually written and any
// errors encountered.
func (s *StandardSyncer) Write(buffer []byte) (int, error) {
	if s.buffer != nil {
		if (len(s.buffer) + len(buffer)) >= s.capacity {
			size, err := s.writer.Write(s.buffer)

			if err != nil {
				s.buffer = s.buffer[size : ]
				return 0, err
			}

			s.buffer = s.buffer[ : 0]
		}

		s.buffer = append(s.buffer, buffer...)
		return len(buffer), nil
	}

	return s.writer.Write(buffer)
}

// Sync writes the internally cached data to a specific storage device.
// If the specific storage device is based on the file system, write the
// data cached by the file system to the persistent storage device.
//
// Finally, any errors encountered are returned.
func (s *StandardSyncer) Sync() error {
	if len(s.buffer) > 0 {
		_, err := s.writer.Write(s.buffer)

		if err != nil {
			return err
		}
	}

	handle, ok := s.writer.(*os.File)

	if !ok {
		return nil
	}

	return handle.Sync()
}

// Close automatically flushes the internal cache once, and then releases
// any kernel objects that have been opened (including but not limited to:
// file handles, etc.).
//
// Finally, any errors encountered are returned.
func (s *StandardSyncer) Close() error {
	s.Sync()
	return nil
}

// StandardSyncerOption is a structure containing standard synchronizer
// options.
type StandardSyncerOption struct {
	SyncerOption

	// Writer represents an instance of a specific storage device that
	// implements io.Writer. If not provided, the default value is
	// ioutil.Discard.
	Writer io.Writer
}

// UseCacheCapacity uses the given capacity as the value of the option
// CacheCapacity. For details, please refer to the comment section of
// the CacheCapacity option. Then return to the option instance itself.
func (o *StandardSyncerOption) UseCacheCapacity(capacity int) *StandardSyncerOption {
	o.CacheCapacity = capacity
	return o
}

// UseWriter uses the given writer as the value of the option Writer.
// If the value of the given writer is nil, ioutil.Discard is used.
// For details, please refer to the comment section of the Writer option.
// Then return to the option instance itself.
func (o *StandardSyncerOption) UseWriter(writer io.Writer) *StandardSyncerOption {
	if writer == nil {
		writer = ioutil.Discard
	}

	o.Writer = writer
	return o
}

// Build builds and returns a standard synchronizer instance.
func (o *StandardSyncerOption) Build() (*StandardSyncer, error) {
	var buffer []byte

	if o.CacheCapacity < 1024 && o.CacheCapacity > 0 {
		o.CacheCapacity = 1024
	}

	if o.CacheCapacity > 0 {
		buffer = make([]byte, 0, o.CacheCapacity)
	}

	return &StandardSyncer {
		writer: o.Writer,
		buffer: buffer,
		capacity: o.CacheCapacity,
	}, nil
}

// NewStandardSyncerOption creates and returns a standard synchronizer
// option instance with default optional values.
func NewStandardSyncerOption() *StandardSyncerOption {
	return &StandardSyncerOption {
		SyncerOption: NewSyncerOption(),
		Writer: ioutil.Discard,
	}
}

// NewStandardSyncer creates and returns a standard synchronizer
// instance using the default optional values.
func NewStandardSyncer() (*StandardSyncer, error) {
	return NewStandardSyncerOption().Build()
}

// FileSyncer is the structure of the file synchronizer instance.
//
// The file synchronizer is based on the standard synchronizer and
// uses a file on the local hard disk as a specific storage device.
//
// Please note that file synchronizers are not thread-safe.
type FileSyncer struct {
	StandardSyncer
}

// Close automatically flushes the internal cache once, and then releases
// any kernel objects that have been opened (including but not limited to:
// file handles, etc.).
//
// Finally, any errors encountered are returned.
func (s *FileSyncer) Close() error {
	s.StandardSyncer.Close()
	return s.writer.(*os.File).Close()
}

// FileSyncerOption is a structure containing file synchronizer options.
type FileSyncerOption struct {
	SyncerOption

	// FileName represents the path name of a file on the local disk used
	// as a specific storage device. If not provided, the default value is
	// os.DevNull.
	FileName string
}

// UseCacheCapacity uses the given capacity as the value of the option
// CacheCapacity. For details, please refer to the comment section of
// the CacheCapacity option. Then return to the option instance itself.
func (o *FileSyncerOption) UseCacheCapacity(capacity int) *FileSyncerOption {
	o.CacheCapacity = capacity
	return o
}

// UseName uses the given name as the value of the option FileName. For
// details, please refer to the comment section of the FileName option.
func (o *FileSyncerOption) UseName(name string) *FileSyncerOption {
	o.FileName = name
	return o
}

// Build builds and returns a file synchronizer instance.
func (o *FileSyncerOption) Build() (*FileSyncer, error) {
	if len(o.FileName) == 0 {
		o.FileName = os.DevNull
	}

	handle, err := os.OpenFile(o.FileName, os.O_RDWR |
		os.O_CREATE | os.O_APPEND, 0666)

	if err != nil {
		return nil, err
	}

	option := NewStandardSyncerOption()

	option.SyncerOption = o.SyncerOption
	option.Writer = handle
	
	syncer, err := option.Build()

	if err != nil {
		handle.Close()
		return nil, err
	}

	return &FileSyncer {
		StandardSyncer: *syncer,
	}, nil
}

// NewFileSyncerOption creates and returns an instance of a file
// synchronizer option with default optional values.
func NewFileSyncerOption() *FileSyncerOption {
	return &FileSyncerOption {
		SyncerOption: NewSyncerOption(),
		FileName: os.DevNull,
	}
}

// NewFileSyncer creates and returns a file synchronizer instance
// using the default optional values.
func NewFileSyncer() (*FileSyncer, error) {
	return NewFileSyncerOption().Build()
}

// DiscardSyncer is the structure of the lost synchronizer instance.
//
// The discard synchronizer is based on the standard synchronizer,
// using ioutil.Discard as an instance of a specific storage device.
// The discard synchronizer will unconditionally discard all written
// log entry data.
//
// Please note that file synchronizers are not thread-safe.
type DiscardSyncer struct {
	StandardSyncer
}

// NewDiscardSyncer creates an instance of a discard synchronizer
// using the default optional values.
func NewDiscardSyncer() (*DiscardSyncer, error) {
	syncer, err := NewStandardSyncerOption().
		UseCacheCapacity(0).Build()

	if err != nil {
		return nil, err
	}

	return &DiscardSyncer {
		StandardSyncer: *syncer,
	}, nil
}
