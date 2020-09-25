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
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
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
	// If not provided, the default value is 32 KB * the number of
	// logical processors.
	CacheCapacity int
}

// NewSyncerOption returns the value of a synchronizer option with the
// default optional value.
func NewSyncerOption() SyncerOption {
	return SyncerOption {
		CacheCapacity: (1024 * 32) * runtime.NumCPU(),
	}
}

// StandardSyncer is the structure of a standard synchronizer instance.
//
// The standard synchronizer uses an instance that implements the io.Writer
// interface as a specific storage device, but it does not check the actual
// data type of the instance when it is closed, and it does not close the
// instance (if supported).
//
// Please note that if the mutex is disabled, the API provided by
// the synchronizer is not thread-safe.
type StandardSyncer struct {
	writer io.Writer
	buffer []byte
	capacity int
	mutex *SpinLock
}

// flush writes the data stored in the internal cache to a specific storage
// device, and then returns the actual number of bytes written and any
// errors encountered.
//
// Please note that this function is not thread-safe.
func (s *StandardSyncer) flush() (int, error) {
	size, err := s.writer.Write(s.buffer)

	if err != nil {
		if size > 0 {
			s.buffer = append(s.buffer[ : 0], s.buffer[size : ]...)
		}

		return size, err
	}

	s.buffer = s.buffer[ : 0]
	return size, nil
}

// Write writes the data of a given buffer slice to a specific storage
// device. If the internal cache is enabled, the internal cache is
// written first. If the capacity of the internal cache is saturated,
// it is automatically flushed once.
//
// Finally, it returns the number of bytes actually written and any
// errors encountered.
func (s *StandardSyncer) Write(buffer []byte) (int, error) {
	if s.mutex != nil {
		s.mutex.Lock()
	}

	if s.buffer != nil {
		if (len(s.buffer) + len(buffer)) >= s.capacity {
			_, err := s.flush()

			if err != nil {
				if s.mutex != nil {
					s.mutex.Unlock()
				}

				return 0, err
			}
		}

		s.buffer = append(s.buffer, buffer...)

		if s.mutex != nil {
			s.mutex.Unlock()
		}

		return len(buffer), nil
	}

	size, err := s.writer.Write(buffer)

	if s.mutex != nil {
		s.mutex.Unlock()
	}

	return size, err
}

// Sync writes the internally cached data to a specific storage device.
// If the specific storage device is based on the file system, write the
// data cached by the file system to the persistent storage device.
//
// Finally, any errors encountered are returned.
func (s *StandardSyncer) Sync() error {
	if s.mutex != nil {
		s.mutex.LockAndSuspend()
	}

	if len(s.buffer) > 0 {
		_, err := s.flush()

		if err != nil {
			if s.mutex != nil {
				s.mutex.UnlockAndResume()
			}

			return err
		}
	}

	handle, ok := s.writer.(*os.File)

	if !ok {
		if s.mutex != nil {
			s.mutex.UnlockAndResume()
		}

		return nil
	}

	err := handle.Sync()

	if s.mutex != nil {
		s.mutex.UnlockAndResume()
	}

	return err
}

// Close automatically flushes the internal cache once, and then releases
// any kernel objects that have been opened (including but not limited to:
// file handles, etc.).
//
// Finally, any errors encountered are returned.
func (s *StandardSyncer) Close() error {
	_ = s.Sync()
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

	// DisableMutex indicates whether to disable the write mutex
	// protection inside the synchronizer. If a particular storage
	// device is thread-safe, disabling the internal write mutex
	// protection can improve performance. If disabled, the internal
	// cache will also be disabled. If not provided, the default value
	// is false.
	DisableMutex bool
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
	var mutex *SpinLock

	if !o.DisableMutex {
		if o.CacheCapacity < 1024 && o.CacheCapacity > 0 {
			o.CacheCapacity = 1024
		}
	
		if o.CacheCapacity > 0 {
			buffer = make([]byte, 0, o.CacheCapacity)
		}

		mutex = NewSpinLock()
	}

	return &StandardSyncer {
		writer: o.Writer,
		buffer: buffer,
		capacity: o.CacheCapacity,
		mutex: mutex,
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
// Please note that if the mutex is disabled, the API provided by
// the synchronizer is not thread-safe.
type FileSyncer struct {
	*StandardSyncer
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
		StandardSyncer: syncer,
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

// NetworkSyncer is the structure of an instance of a network
// synchronizer.
//
// The network synchronizer is based on the standard synchronizer
// and uses TCP/IP or Unix streams as a specific storage device.
//
// Please note that if the mutex is disabled, the API provided by
// the synchronizer is not thread-safe.
type NetworkSyncer struct {
	*StandardSyncer

	protocol string
	address string

	context context.Context
	contextCancel context.CancelFunc
	contextWaitGroup *sync.WaitGroup

	disconnected int32
}

func (s *NetworkSyncer) reconnect() {
	defer s.contextWaitGroup.Done()

	dialer := &net.Dialer {
		Timeout: time.Second * 5,
	}

	for {
		connect, err := dialer.DialContext(s.context, s.protocol, s.address)

		if err != nil {
			// If the synchronizer is closing, give up the reconnection
			// and return. To avoid calling the function again, the value
			// of `s.disconnected` is not reset.
			if s.context.Err() != nil {
				return
			}

			// Reconnection failed, try again after an interval of 1
			// second.
			select {
			case <-time.After(1 * time.Second):
				continue
			case <-s.context.Done():
				return
			}
			
		}

		s.writer.(net.Conn).Close()
		s.writer = connect
		break
	}

	atomic.CompareAndSwapInt32(&s.disconnected, 1, 0)
}

// Write writes the data of a given buffer slice to a specific storage
// device. If the internal cache is enabled, the internal cache is
// written first. If the capacity of the internal cache is saturated,
// it is automatically flushed once.
//
// Finally, it returns the number of bytes actually written and any
// errors encountered.
func (s *NetworkSyncer) Write(buffer []byte) (int, error) {
	size, err := s.StandardSyncer.Write(buffer)

	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			// The connection to the other end of the network may have
			// been interrupted unexpectedly, try to re-establish the
			// connection.
			if atomic.CompareAndSwapInt32(&s.disconnected, 0, 1) {
				s.contextWaitGroup.Add(1)
				go s.reconnect()
			}
		}
	}

	return size, err
}

// Close automatically flushes the internal cache once, and then releases
// any kernel objects that have been opened (including but not limited to:
// network handles, etc.).
//
// Finally, any errors encountered are returned.
func (s *NetworkSyncer) Close() error {
	s.contextCancel()
	s.contextWaitGroup.Wait()
	s.StandardSyncer.Close()
	return s.StandardSyncer.writer.(net.Conn).Close()
}

const (
	// ProtocolTCP represents that the communication protocol of the
	// network synchronizer is TCP/IP. For details, please refer to the
	// comment section of the NetworkSyncer structure.
	ProtocolTCP = "tcp"

	// ProtocolUnix represents that the communication protocol of the
	// network synchronizer is Unix Domain Socket. For details, please
	// refer to the comment section of the NetworkSyncer structure.
	ProtocolUnix = "unix"
)

var (
	// ErrInvalidProtocol represents that the network protocol is invalid
	// or unsupported. This is usually because the value of the given
	// network protocol type is invalid.
	ErrInvalidProtocol = errors.New("invalid network protocol")
)

// NetworkSyncerOption is a structure containing network synchronizer
// options.
type NetworkSyncerOption struct {
	SyncerOption

	// Protocol represents the communication protocol used by the network
	// synchronizer, which will determine how the network synchronizer
	// establishes a connection with the other end of the network. The
	// optional values are defined by the constants at the beginning of
	// Protocol...
	//
	// If not provided, the default value is the ProtocolUnix constant.
	Protocol string
	
	// Address represents the address of the other end of the network
	// where the network synchronizer uses a specific communication
	// protocol to establish a connection. The address format depends on
	// the value of the Protocol option.
	//
	// If not provided, the default value is /var/run/santa.sock. It is
	// worth noting that the default value is invalid for Windows.
	Address string
}

// UseCacheCapacity uses the given capacity as the value of the option
// CacheCapacity. For details, please refer to the comment section of
// the CacheCapacity option. Then return to the option instance itself.
func (o *NetworkSyncerOption) UseCacheCapacity(capacity int) *NetworkSyncerOption {
	o.CacheCapacity = capacity
	return o
}

// UseProtocol uses the given protocol as the value of the option Protocol.
// Please refer to the comment section of the Protocol option for details.
// Then return to the option instance itself.
func (o *NetworkSyncerOption) UseProtocol(protocol string) *NetworkSyncerOption {
	o.Protocol = protocol
	return o
}

// UseAddress uses the given address as the value of the option Address,
// please refer to the comment section of the Address option for details.
// Then return to the option instance itself.
func (o *NetworkSyncerOption) UseAddress(address string) *NetworkSyncerOption {
	o.Address = address
	return o
}

// Build builds and returns an instance of the network synchronizer and
// any errors encountered.
func (o *NetworkSyncerOption) Build() (*NetworkSyncer, error) {
	switch o.Protocol {
	case ProtocolTCP:
	case ProtocolUnix:
	default:
		return nil, ErrInvalidProtocol
	}

	connect, err := net.Dial(o.Protocol, o.Address)
	
	if err != nil {
		return nil, err
	}

	option := NewStandardSyncerOption()

	option.SyncerOption = o.SyncerOption
	option.Writer = connect
	
	syncer, err := option.Build()

	if err != nil {
		connect.Close()
		return nil, err
	}

	context, contextCancel := context.WithCancel(
		context.Background())

	return &NetworkSyncer {
		StandardSyncer: syncer,

		protocol: o.Protocol,
		address: o.Address,

		context: context,
		contextCancel: contextCancel,
		contextWaitGroup: &sync.WaitGroup { },
	}, nil
}

// NewNetworkSyncerOption creates and returns a network synchronizer
// option instance with default option values.
func NewNetworkSyncerOption() *NetworkSyncerOption {
	return &NetworkSyncerOption {
		SyncerOption: NewSyncerOption(),
		Protocol: ProtocolUnix,
		Address: "/var/run/santa.sock",
	}
}

// DiscardSyncer is the structure of the discard synchronizer instance.
//
// The discard synchronizer is based on the standard synchronizer,
// using ioutil.Discard as an instance of a specific storage device.
// The discard synchronizer will unconditionally discard all written
// log entry data.
//
// Please note that file synchronizers are not thread-safe.
type DiscardSyncer struct {
	*StandardSyncer
}

// NewDiscardSyncer creates an instance of a discard synchronizer
// using the default optional values.
func NewDiscardSyncer() (*DiscardSyncer, error) {
	option := NewStandardSyncerOption()
	option.Writer = ioutil.Discard
	option.DisableMutex = true

	syncer, err := option.Build()

	if err != nil {
		return nil, err
	}

	return &DiscardSyncer {
		StandardSyncer: syncer,
	}, nil
}
