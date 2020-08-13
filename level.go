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
	"errors"
	"fmt"
	"strings"
)

// Level is a data type represents the log level.
type Level uint8

const (
	// LevelDebug means the log level DEBUG, usually used to record
	// development and debugging logs.
	LevelDebug Level = iota
	
	// LevelInfo represents the log level INFO, usually used to record
	// regular logs.
	LevelInfo

	// LevelWarning represents the log level WARNING, which is usually
	// used to record normal but important logs.
	LevelWarning

	// LevelError means log level ERROR, usually used to record errors
	// but not fatal logs.
	LevelError

	// LevelFatal represents the log level FATAL, usually used to record
	// fatal error logs.
	LevelFatal
)

var (
	// ErrorLevelInvalid represents the given invalid log level.
	ErrorLevelInvalid = errors.New("invalid level")
)

// Enabled checks whether the given log level is enabled.
func (l Level) Enabled(level Level) bool {
	return l <= level
}

// String Returns the name string of the log level.
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarning:
		return "warning"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return fmt.Sprintf("unknown(%d)", l)
	}
}

// Format returns the formatting style string of the log level.
func (l Level) Format() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarning:
		return "WARNING"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", l)
	}
}

// AppendFormat appends the format string of the log level to the
// given buffer slice, and then returns the appended buffer slice.
func (l Level) AppendFormat(buffer []byte) []byte {
	switch l {
	case LevelDebug:
		return append(buffer, "DEBUG"...)
	case LevelInfo:
		return append(buffer, "INFO"...)
	case LevelWarning:
		return append(buffer, "WARNING"...)
	case LevelError:
		return append(buffer, "ERROR"...)
	case LevelFatal:
		return append(buffer, "FATAL"...)
	default:
		return append(buffer, fmt.Sprintf("UNKNOWN(%d)", l)...)
	}
}

// ParseFrom parse the log level from the name.
func (l Level) ParseFrom(name string) error {
	switch strings.ToLower(name) {
	case "debug":
		l = LevelDebug
	case "info":
		l = LevelInfo
	case "warning":
		l = LevelWarning
	case "error":
		l = LevelError
	case "fatal":
		l = LevelFatal
	default:
		return ErrorLevelInvalid
	}

	return nil
}

// ParseLevel parse the log level from the name.
func ParseLevel(name string) (Level, error) {
	var level Level

	err := level.ParseFrom(name)

	if err != nil {
		return 0, err
	}

	return level, nil
}

// LevelSpan is a structure that contains the log level span.
type LevelSpan struct {
	// Start represents the starting level of the log.
	Start Level

	// End represents the end level of the log level.
	End Level
}

// Contains checks whether the given log level is within the span.
func (l LevelSpan) Contains(level Level) bool {
	return level >= l.Start && level <= l.End
}
