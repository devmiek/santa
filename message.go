// message.go is the golang-1.14.6 source file.

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
	Fields ElementFields
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
	buffer = append(buffer, "\", \"payload\": "...)
	return m.Fields.FormatJSON(buffer)
}

// SampleText returns the text sample string of the log entry message.
func (m StructMessage) SampleText() string {
	return m.Text
}
