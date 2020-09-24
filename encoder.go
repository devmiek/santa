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
	"strconv"
	"time"
)

var (
	// ErrUnsupportedMessage represents that the message type of the
	// given log entry is not supported, usually because the message
	// does not correctly implement the message serialization interface
	// of the encoder.
	ErrUnsupportedMessage = errors.New("unsupported message type")
)

// EncoderOption is a structure that contains options for the encoder.
//
// Encoder options include basic options for all types of encoder options.
// Normally, any encoder option type should include this structure.
type EncoderOption struct {
	// EncodeTime represents whether to encode the time of the log entry
	// and append it to the encoding result. If not provided, the default
	// value is true.
	EncodeTime bool

	// EncodeSourceLocation represents whether the output function call
	// source location of the log entry is encoded and appended to the
	// encoding result. If not provided, the default value is true.
	EncodeSourceLocation bool

	// EncodeLabels represents whether to encode one or more labels of
	// the log entry and append them to the encoding result. If not
	// provided, the default value is true.
	EncodeLabels bool
	
	// EncodeName represents whether to encode the name of the log entry
	// and append it to the encoding result. If not provided, the default
	// value is true.
	EncodeName bool

	// EncodeLevel represents whether to encode the level of the log entry
	// and append it to the encoding result. If not provided, the default
	// value is true.
	EncodeLevel bool
}

// NewEncoderOption returns an encoder option value with default optional
// values.
func NewEncoderOption() EncoderOption {
	return EncoderOption {
		EncodeTime: true,
		EncodeSourceLocation: true,
		EncodeLabels: true,
		EncodeName: true,
		EncodeLevel: true,
	}
}

// Encoder is the public interface of the encoder.
//
// The encoder encodes log entries into consecutive bytes in a specific
// format. For example, the standard encoder encodes log entries into
// multi-field strings that can be easily read by humans; the JSON encoder
// encodes log entries into JSON strings that can be easily parsed by
// machines.
//
// The message type of any log entry must implement the serializer
// interface provided by the specific encoder, otherwise the encoder does
// not know how to encode the message.
type Encoder interface {
	// Encode encodes a given log entry into consecutive bytes in a specific
	// format, then appends to the given buffer slice, and finally returns
	// the appended buffer slice.
	Encode(buffer []byte, entry *Entry) ([]byte, error)

	// Option returns the value of the basic options of the encoder, and the
	// application can optimize the actual behavior by checking the values
	// of the options.
	Option() EncoderOption
}

// EncoderKeys is a structure containing the key names used when encoding
// log entries.
//
// This structure contains the name of the key used when encoding each
// structured log entry. It is usually used to customize the key of the
// structured log after encoding. For example, the log level key defaults
// to "level", the log message key defaults to "message", etc. .
//
// Please note that the encoder does not check the key name, which means
// that the key name string is allowed to be empty, but this is not in
// compliance with the specification for some structured encoders.
type EncoderKeys struct {
	// TimeKey represents the name of the key used when encoding the time
	// of the log entry. If not provided, the default value is "timestamp".
	TimeKey string

	// SourceLocationKey represents the name of the key used when the output
	// of the encoded log entry calls the source location. If not provided,
	// the default value is "sourceLocation".
	SourceLocationKey string

	// LabelsKey represents the name of the key used when encoding a set of
	// labels for log entries. If not provided, the default value is "labels".
	LabelsKey string

	// NameKey represents the name of the key used when encoding the name of
	// the log entry. If not provided, the default value is "logName".
	NameKey string

	// LevelKey represents the name of the key used when encoding the level
	// of a log entry. If not provided, the default value is "level".
	LevelKey string

	// MessageKey represents the name of the key used when encoding the
	// message of the log entry. If not provided, the default value is
	// "message".
	MessageKey string
}

// NewEncoderKeys returns an EncoderKeys value with the name of the key
// of the default log entry.
func NewEncoderKeys() EncoderKeys {
	return EncoderKeys {
		TimeKey: "timestamp",
		SourceLocationKey: "sourceLocation",
		LabelsKey: "labels",
		NameKey: "logName",
		LevelKey: "level",
		MessageKey: "message",
	}
}

// StandardSerializer is the public interface of the standard serializer.
//
// Any message type of a log entry encoded by a standard encoder must
// implement this interface, otherwise the standard encoder does not know
// how to encode the message part of a log entry.
type StandardSerializer interface {
	// SerializeStandard serializes the message or any content and appends
	// to the given buffer slice, and then returns the appended buffer
	// slice.
	SerializeStandard(buffer []byte) []byte
}

// StandardEncoder is the structure of a standard encoder instance.
// 
// Standard encoders encode log entries into human-readable strings,
// and are usually used to print log entries on the console. If the
// application does not need to print out structured logs, a standard
// encoder is a good choice.
//
// Please note that the message type of any log entry encoded with a
// standard encoder must implement the StandardSerializer interface,
// otherwise the standard encoder does not know how to encode the
// message part of the log entry.
type StandardEncoder struct {
	layout string
	option EncoderOption
}

// Encode encodes a given log entry into consecutive bytes in a specific
// format, then appends to the given buffer slice, and finally returns
// the appended buffer slice.
func (e *StandardEncoder) Encode(buffer []byte, entry *Entry) ([]byte, error) {
	if e.option.EncodeTime {
		if len(e.layout) == 0 {
			buffer = strconv.AppendInt(buffer, entry.Time.UnixNano(), 10)
		} else {
			buffer = entry.Time.AppendFormat(buffer, e.layout)
		}

		buffer = append(buffer, ' ')
	}

	if e.option.EncodeSourceLocation {
		buffer = entry.SourceLocation.AppendString(buffer)
		buffer = append(buffer, ' ')
	}

	if e.option.EncodeLabels {
		if entry.Labels.Count() == 0 {
			buffer = append(buffer, "no-labels "...)
		} else {
			buffer = entry.Labels.SerializeStandard(buffer)
			buffer = append(buffer, ' ')
		}
	}

	if e.option.EncodeName && len(entry.Name) > 0 {
		buffer = append(buffer, entry.Name...)
		buffer = append(buffer, ' ')
	}

	if e.option.EncodeLevel {
		buffer = append(buffer, '[')
		buffer = append(buffer, entry.Level.Format()...)
		buffer = append(buffer, "] "...)
	}

	switch message := entry.Message.(type) {
	case nil:
		buffer = append(buffer, "null"...)
	case StandardSerializer:
		buffer = message.SerializeStandard(buffer)
	default:
		return nil, ErrUnsupportedMessage
	}

	return append(buffer, '\n'), nil
}

// Option returns the value of the basic options of the encoder, and the
// application can optimize the actual behavior by checking the values
// of the options.
func (e *StandardEncoder) Option() EncoderOption {
	return e.option
}

// StandardEncoderOption is a structure that contains options for standard
// encoders.
type StandardEncoderOption struct {
	EncoderOption

	// TimeLayout represents the time formatting layout style used when
	// encoding the time of the log entry. If not provided, the default
	// value depends on the encoder type.
	//
	// If the value of this option is an empty string, the UNIX nanosecond
	// timestamp layout style is used by default.
	TimeLayout string
}

// UseEncoderOption uses the given encoder option as part of the standard
// encoder option. For details, please refer to the comment section of
// the EncoderOption structure. Then return to the option instance itself.
func (o *StandardEncoderOption) UseEncoderOption(option EncoderOption) *StandardEncoderOption {
	o.EncoderOption = option
	return o
}

// UseTimeLayout uses the given layout as the value of the option TimeLayout.
// For details, please refer to the comment section of the TimeLayout option.
// Then return to the option instance itself.
func (o *StandardEncoderOption) UseTimeLayout(layout string) *StandardEncoderOption {
	o.TimeLayout = layout
	return o
}

// Build builds and returns a standard encoder instance.
func (o *StandardEncoderOption) Build() (*StandardEncoder, error) {
	return &StandardEncoder {
		layout: o.TimeLayout,
		option: o.EncoderOption,
	}, nil
}

// NewStandardEncoderOption creates and returns a standard encoder option
// instance with default optional values.
func NewStandardEncoderOption() *StandardEncoderOption {
	return &StandardEncoderOption {
		EncoderOption: NewEncoderOption(),
		TimeLayout: time.RFC3339Nano,
	}
}

// NewStandardEncoder creates and returns a standard encoder instance
// using the default optional values.
func NewStandardEncoder() (*StandardEncoder, error) {
	return NewStandardEncoderOption().Build()
}

// JSONSerializer is the public interface of JSON serializer.
//
// Any message type of a log entry encoded with a JSON encoder must
// implement this interface, otherwise the JSON encoder does not know
// how to encode the message part of a log entry.
type JSONSerializer interface {
	// SerializeJSON serializes the message or any content and appends
	// to the given buffer slice, and then returns the appended buffer
	// slice.
	SerializeJSON(buffer []byte) []byte
}

// JSONEncoder is the structure of the JSON encoder instance.
//
// The JSON encoder is a structured log encoder. The structured
// log will reduce the complexity of the machine to analyze the log.
// It is usually used in a production environment. For the log output
// to the console, it is recommended to choose a standard encoder that
// is easier for humans to read.
//
// Please note that the message type of any log entry encoded with the
// JSON encoder must implement the JSONSerializer interface, otherwise
// the JSON encoder does not know how to encode the message part of the
// log entry.
type JSONEncoder struct {
	layout string
	keys EncoderKeys
	option EncoderOption
}

// Encode encodes a given log entry into consecutive bytes in a specific
// format, then appends to the given buffer slice, and finally returns
// the appended buffer slice.
func (e *JSONEncoder) Encode(buffer []byte, entry *Entry) ([]byte, error) {
	message, ok := entry.Message.(JSONSerializer)

	if !ok {
		return nil, ErrUnsupportedMessage
	}

	buffer = append(buffer, '{')

	if e.option.EncodeTime {
		buffer = append(buffer, '"')
		buffer = append(buffer, e.keys.TimeKey...)

		if len(e.layout) == 0 {
			buffer = append(buffer, "\": "...)
			buffer = strconv.AppendInt(buffer, entry.Time.UnixNano(), 10)
			buffer = append(buffer, ", "...)
		} else {
			buffer = append(buffer, "\": \""...)
			buffer = entry.Time.AppendFormat(buffer, e.layout)
			buffer = append(buffer, "\", "...)
		}
	}

	if e.option.EncodeSourceLocation {
		buffer = append(buffer, '"')
		buffer = append(buffer, e.keys.SourceLocationKey...)
		buffer = append(buffer, "\": "...)
		buffer = entry.SourceLocation.SerializeJSON(buffer)
		buffer = append(buffer, ", "...)
	}

	if e.option.EncodeLabels {
		buffer = append(buffer, '"')
		buffer = append(buffer, e.keys.LabelsKey...)
		buffer = append(buffer, `": `...)

		if entry.Labels.Count() == 0 {
			buffer = append(buffer, "null"...)
		} else {
			buffer = entry.Labels.SerializeJSON(buffer)
		}

		buffer = append(buffer, ", "...)
	}

	if e.option.EncodeName {
		buffer = append(buffer, '"')
		buffer = append(buffer, e.keys.NameKey...)

		if len(entry.Name) > 0 {
			buffer = append(buffer, "\": \""...)
			buffer = append(buffer, entry.Name...)
			buffer = append(buffer, "\", "...)
		} else {
			buffer = append(buffer, "\": "...)
			buffer = append(buffer, "null"...)
			buffer = append(buffer, ", "...)
		}
	}

	if e.option.EncodeLevel {
		buffer = append(buffer, '"')
		buffer = append(buffer, e.keys.LevelKey...)
		buffer = append(buffer, "\": \""...)
		buffer = entry.Level.AppendFormat(buffer)
		buffer = append(buffer, "\", "...)
	}

	buffer = append(buffer, '"')
	buffer = append(buffer, e.keys.MessageKey...)
	buffer = append(buffer, "\": "...)
	buffer = message.SerializeJSON(buffer)

	return append(buffer, "}\n"...), nil
}

// Option returns the value of the basic options of the encoder, and the
// application can optimize the actual behavior by checking the values
// of the options.
func (e *JSONEncoder) Option() EncoderOption {
	return e.option
}

// JSONEncoderOption is a structure containing options for the JSON encoder.
type JSONEncoderOption struct {
	StandardEncoderOption
	EncoderKeys
}

// UseEncoderOption uses the given encoder option as part of the JSON
// encoder option. For details, please refer to the comment section of
// the EncoderOption structure. Then return to the option instance itself.
func (o *JSONEncoderOption) UseEncoderOption(option EncoderOption) *JSONEncoderOption {
	o.EncoderOption = option
	return o
}

// UseEncoderKeys uses the given encoder keys as part of the JSON encoder
// options. For details, please refer to the comments section of the
// EncoderKeys structure. Then return to the option instance itself.
func (o *JSONEncoderOption) UseEncoderKeys(keys EncoderKeys) *JSONEncoderOption {
	o.EncoderKeys = keys
	return o
}

// Build builds and returns an instance of the JSON encoder.
func (o *JSONEncoderOption) Build() (*JSONEncoder, error) {
	return &JSONEncoder {
		layout: o.TimeLayout,
		keys: o.EncoderKeys,
		option: o.EncoderOption,
	}, nil
}

// NewJSONEncoderOption creates and returns a JSON encoder option instance
// with default optional values.
func NewJSONEncoderOption() *JSONEncoderOption {
	return &JSONEncoderOption {
		StandardEncoderOption: *NewStandardEncoderOption().UseTimeLayout(""),
		EncoderKeys: NewEncoderKeys(),
	}
}

// NewJSONEncoder creates and returns a standard encoder instance
// using the default optional values.
func NewJSONEncoder() (*JSONEncoder, error) {
	return NewJSONEncoderOption().Build()
}
