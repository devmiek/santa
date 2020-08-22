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
	"fmt"
)

// Message is the public interface for messages.
type Message interface { }

// StringMessage is the data type of the string log entry message.
type StringMessage string

// FormatStandard formats the log entry message as a string and appends
// it to the given buffer slice, and then returns the appended buffer slice.
func (m StringMessage) FormatStandard(buffer []byte) []byte {
	return append(buffer, m...)
}

// FormatJSON formats the log entry message into a JSON string and appends
// it to the given buffer slice, and then returns the appended buffer slice.
func (m StringMessage) FormatJSON(buffer []byte) []byte {
	buffer = append(buffer, '"')
	buffer = append(buffer, m...)
	return append(buffer, '"')
}

// SampleText returns the text sample string of the log entry message.
func (m StringMessage) SampleText() string {
	return string(m)
}

// TemplateMessage is a message structure containing formatted
// templates and parameter values.
type TemplateMessage struct {
	// Template represents the template string of the message,
	// and the message encoder determines the string format of
	// the message through the template string.
	Template string
	
	// Args represents the formatting parameters of the template
	// message. The number and position of the parameters correspond
	// to the template string.
	Args []interface { }
}

// FormatStandard formats the log entry message as a string and appends
// it to the given buffer slice, and then returns the appended buffer slice.
func (m TemplateMessage) FormatStandard(buffer []byte) []byte {
	return append(buffer, fmt.Sprintf(m.Template, m.Args...)...)
}

// FormatJSON formats the log entry message into a JSON string and appends
// it to the given buffer slice, and then returns the appended buffer slice.
func (m TemplateMessage) FormatJSON(buffer []byte) []byte {
	buffer = append(buffer, '"')
	buffer = append(buffer, fmt.Sprintf(m.Template, m.Args...)...)
	return append(buffer, '"')
}

// SampleText returns the text sample string of the log entry message.
func (m TemplateMessage) SampleText() string {
	return m.Template
}

// StructMessage is a log entry message structure containing
// multiple fields.
type StructMessage struct {
	// Text represents the message text, usually the message
	// text is used to summarize the subject of the log entry.
	Text string

	// Fields represents one or more fields included in the
	// field message, and these fields will be encoded as
	// structured log entries.
	Fields ElementObject
}

// FormatStandard formats the log entry message as a string and appends
// it to the given buffer slice, and then returns the appended buffer slice.
func (m StructMessage) FormatStandard(buffer []byte) []byte {
	buffer = append(buffer, m.Text...)
	buffer = append(buffer, ' ')
	return m.Fields.FormatJSON(buffer)
}

// FormatJSON formats the log entry message into a JSON string and appends
// it to the given buffer slice, and then returns the appended buffer slice.
func (m StructMessage) FormatJSON(buffer []byte) []byte {
	buffer = append(buffer, '"')
	buffer = append(buffer, m.Text...)

	if len(m.Fields) == 0 {
		return append(buffer, '"')
	}

	buffer = append(buffer, "\", \"payload\": "...)
	return m.Fields.FormatJSON(buffer)
}

// SampleText returns the text sample string of the log entry message.
func (m StructMessage) SampleText() string {
	return m.Text
}
