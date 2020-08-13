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
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

// EntrySourceLocation is a structure that contains the source of log entries.
type EntrySourceLocation struct {
	// Proc represents the address of the caller function that
	// printed the log entry.
	Proc uintptr
	
	// File represents the file path of the source code of the
	// caller function for printing log entries.
	File string
	
	// Line represents the line number of the source file where
	// the caller of the log entry is printed.
	Line int
	
	// Parsed represents whether the source of the log entry has
	// been successfully parsed.
	Parsed bool
}

// AppendString encodes the source of the log entry as a string, then
// appends it to the end of the given buffer slice, and finally
// returns the new buffer slice.
func (s *EntrySourceLocation) AppendString(buffer []byte) []byte {
	if buffer == nil {
		return nil
	}

	if !s.Parsed {
		return append(buffer, "???:0"...)
	}

	buffer = append(buffer, filepath.Base(s.File)...)
	buffer = append(buffer, ':')

	return strconv.AppendInt(buffer, int64(s.Line), 10)
}

// AppendJSON encodes the source location of the log entry as a JSON
// string and appends it to the given buffer slice, and then returns
// the appended buffer slice.
func (s *EntrySourceLocation) AppendJSON(buffer []byte) []byte {
	if buffer == nil {
		return nil
	}

	if !s.Parsed {
		return append(buffer, "null"...)
	}

	buffer = append(buffer, "{\"file\": \""...)
	buffer = append(buffer, filepath.Base(s.File)...)
	buffer = append(buffer, "\", \"line\": "...)
	buffer = strconv.AppendInt(buffer, int64(s.Line), 10)
	buffer = append(buffer, ", \"function\": \""...)
	buffer = append(buffer, runtime.FuncForPC(s.Proc).Name()...)

	return append(buffer, "\"}"...)
}

// newEntrySourceLocation receives the return value of the runtime.Caller
// function to facilitate the creation of the value of the source location
// of the log entry.
func newEntrySourceLocation(p uintptr, f string, l int, ok bool) EntrySourceLocation {
	return EntrySourceLocation {
		Proc: p,
		File: f,
		Line: l,
		Parsed: ok,
	}
}

// Entry is the structure of the log entry instance.
type Entry struct {
	// Time represents the generation time of the log entry, usually
	// the time when the log entry is printed out.
	Time time.Time

	// Level represents the severity level of the log entry.
	Level Level

	// Message represents the message instance of the log entry. The
	// message instance will be encoded into the target data format by
	// the log entry encoder.
	//
	// The value can be nil.
	Message Message

	// SourceLocation represents the source code location of the log
	// entry, usually the calling location of the log entry operation.
	SourceLocation EntrySourceLocation

	// Name represents the name of the log entry.
	Name string
}
