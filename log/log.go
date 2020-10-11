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

package log

import (
	"errors"

	"github.com/nobody-night/santa"
)

var (
	// pool is a structural variable that contains default instances of
	// various pools. These pool instances are automatically created when
	// the application is initialized and shared globally.
	pool santa.GlobalPool = santa.GetGlobalPool()

	// logger is an instance of the standard logger, which is used as the
	// default logger instance. The default logger instance is automatically
	// created when the application is initialized and shared globally.
	logger *santa.StandardLogger = nil
)

// init initializes the default standard logger instance.
func init() {
	option := santa.NewStandardOption()
	option.Encoding.UseStandard()
	instance, err := option.Build()
	if err != nil {
		panic(err)
	}
	logger = instance
}

// Set sets the default logger to the given standard logger instance and
// returns any errors encountered. This function will try to close the old
// default logger instance.
//
// Please note that this API is not thread-safe.
func Set(instance *santa.StandardLogger) error {
	err := logger.Close()
	if err != nil && !errors.Is(err, santa.ErrClosed) {
		return err		
	}
	logger = instance
	return nil
}

// Close close all specific exporters, and then return any errors
// encountered. For details, please refer to the comment section of the
// Close function of the Exporter interface.
//
// If there are multiple copies of the logger, this function only reduces
// the reference count of the logger. If the logger's reference count is 0,
// it will actually be closed.
//
// Please note that this function is not guaranteed to succeed. If any
// errors are encountered, the state of the application may change. The
// best practice is to exit the application.
func Close() error {
	return logger.Close()
}

// Sync writes the internal cache data of a specific synchronizer to a
// specific storage device. If the specific storage device is based on
// the file system, write the data cached by the file system to the
// persistent storage device. For details, please refer to the Sync
// function of the Syncer interface.
//
// Finally, any errors encountered are returned.
func Sync() error {
	return logger.Sync()
}

// Duplicate creates and returns a copy of the logger. If the logger is
// closed, it returns nil.
//
// Please note that the application must explicitly close each copy of
// the logger, otherwise the logger may be leaked.
func Duplicate() *santa.StandardLogger {
	return logger.Duplicate()
}

// SetName sets the log entry name to the given name. For details, please
// refer to the comment section of the Name field of the StandardOption
// structure.
//
// Please note that this API is not thread-safe.
func SetName(name string) {
	logger.SetName(name)
}

// SetLevel sets the lowest level of the log entry to the given level.
// For details, please refer to the comment section of the Level field of
// the StandardOption structure.
//
// Please note that this API is not thread-safe.
func SetLevel(level santa.Level) {
	logger.SetLevel(level)
}

// SetSampler sets the sampler to the given sampler. For details, please
// refer to the comment section of the Sampler field of the Option
// structure.
//
// Please note that this API is not thread-safe.
func SetSampler(sampler santa.Sampler) {
	logger.SetSampler(sampler)
}

// SetLabels sets the label to one or more given labels. For details,
// please refer to the comment section of the Labels field of the Option
// structure.
//
// It is worth noting that one or more labels previously set by the
// logger will be discarded because labels need to be pre-serialized.
//
// Please note that this API is not thread-safe.
func SetLabels(labels ...santa.Label) {
	logger.SetLabels(labels...)
}

// AddHooks adds one or more hooks to the hook chain. For details,
// please refer to the comment section of the Hooks field of the Option
// option.
//
// Please note that this API is not thread-safe.
func AddHooks(hooks ...santa.Hook) {
	logger.AddHooks(hooks...)
}

// ResetHooks resets the hook chain, and the hooks that have been added
// will be removed. For details, please refer to the comment section of
// the Hooks field of the Option option.
//
// Please note that this API is not thread-safe.
func ResetHooks() {
	logger.ResetHooks()
}

// Prints outputs a structured log message with a given log level,
// given description text and fields, and then returns any errors
// encountered.
func Prints(level santa.Level, text string, fields ...santa.Field) error {
	message := pool.Message.Structure.New(text, fields)
	err := logger.Output(2, level, message)
	pool.Message.Structure.Free(message)
	return err
}

// Debugs outputs a structured log message with a log level of DEBUG,
// given description text and fields, and then returns any errors
// encountered.
func Debugs(text string, fields ...santa.Field) error {
	message := pool.Message.Structure.New(text, fields)
	err := logger.Output(2, santa.LevelDebug, message)
	pool.Message.Structure.Free(message)
	return err
}

// Infos outputs a structured log message with a log level of INFO,
// given description text and fields, and then returns any errors
// encountered.
func Infos(text string, fields ...santa.Field) error {
	message := pool.Message.Structure.New(text, fields)
	err := logger.Output(2, santa.LevelInfo, message)
	pool.Message.Structure.Free(message)
	return err
}

// Warnings outputs a structured log message with a log level of WARNING,
// given description text and fields, and then returns any errors
// encountered.
func Warnings(text string, fields ...santa.Field) error {
	message := pool.Message.Structure.New(text, fields)
	err := logger.Output(2, santa.LevelWarning, message)
	pool.Message.Structure.Free(message)
	return err
}

// Errors outputs a structured log message with a log level of ERROR,
// given description text and fields, and then returns any errors
// encountered.
func Errors(text string, fields ...santa.Field) error {
	message := pool.Message.Structure.New(text, fields)
	err := logger.Output(2, santa.LevelError, message)
	pool.Message.Structure.Free(message)
	return err
}

// Fatals outputs a structured log message with a log level of FATAL,
// given description text and fields, and then returns any errors
// encountered.
func Fatals(text string, fields ...santa.Field) error {
	message := pool.Message.Structure.New(text, fields)
	err := logger.Output(2, santa.LevelFatal, message)
	pool.Message.Structure.Free(message)
	return err
}

// Printf outputs a template log message with a given log level, a given
// template string and one or more parameters, and then returns any errors
// encountered.
func Printf(level santa.Level, template string, args ...interface { }) error {
	message := pool.Message.Template.New(template, args)
	err := logger.Output(2, level, message)
	pool.Message.Template.Free(message)
	return err
}

// Debugf outputs a template log message with a log level of DEBUG, a given
// template string and one or more parameters, and then returns any errors
// encountered.
func Debugf(template string, args ...interface { }) error {
	message := pool.Message.Template.New(template, args)
	err := logger.Output(2, santa.LevelDebug, message)
	pool.Message.Template.Free(message)
	return err
}

// Infof outputs a template log message with a log level of INFO, a given
// template string and one or more parameters, and then returns any errors
// encountered.
func Infof(template string, args ...interface { }) error {
	message := pool.Message.Template.New(template, args)
	err := logger.Output(2, santa.LevelInfo, message)
	pool.Message.Template.Free(message)
	return err
}

// Warningf outputs a template log message with a log level of WARNING, a
// given template string and one or more parameters, and then returns any
// errors encountered.
func Warningf(template string, args ...interface { }) error {
	message := pool.Message.Template.New(template, args)
	err := logger.Output(2, santa.LevelWarning, message)
	pool.Message.Template.Free(message)
	return err
}

// Errorf outputs a template log message with a log level of ERROR, a given
// template string and one or more parameters, and then returns any errors
// encountered.
func Errorf(template string, args ...interface { }) error {
	message := pool.Message.Template.New(template, args)
	err := logger.Output(2, santa.LevelError, message)
	pool.Message.Template.Free(message)
	return err
}

// Fatalf outputs a template log message with a log level of FATAL, a given
// template string and one or more parameters, and then returns any errors
// encountered.
func Fatalf(template string, args ...interface { }) error {
	message := pool.Message.Template.New(template, args)
	err := logger.Output(2, santa.LevelFatal, message)
	pool.Message.Template.Free(message)
	return err
}
