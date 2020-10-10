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

// Label is a structure that contains the label name and value.
//
// Label is a pair of custom static key-value data used to identify the
// same attributes of one or more log entries to provide more information
// about the log entries.
//
// Unlike Field, Label is usually specified when creating the Logger
// instance. For example, use Label to identify that the log entries
// output by one or more loggers belong to a specific system, module, and
// user.
//
// Please note that the API provided by Label is not thread-safe.
type Label struct {
	// Key represents the unique key name of the label, such as "system",
	// "module", or "user".
	Key string
	
	// Value represents the key value of label and can be any string
	// related to the label key name. For example, the name of the system,
	// the name of the module, the unique identifier of the user, etc.
	Value string
}

// SerializeJSON serializes the label into a JSON string and appends it to
// the given buffer slice, and then returns the appended buffer slice.
func (l Label) SerializeJSON(buffer []byte) []byte {
	buffer = append(buffer, '"')
	buffer = append(buffer, l.Key...)
	buffer = append(buffer, `": "`...)
	buffer = append(buffer, l.Value...)
	return append(buffer, '"')
}

// SerializeStandard serializes the label into a standard log string and
// appends it to the given buffer slice, and then returns the appended
// buffer slice.
func (l Label) SerializeStandard(buffer []byte) []byte {
	return l.SerializeJSON(buffer)
}

// NewLabel returns a label value with a given key value.
func NewLabel(key, value string) Label {
	return Label {
		Key: key,
		Value: value,
	}
}

// Labels is a structure containing one or more labels. For details,
// please refer to the annotation section of the Label structure.
type Labels []Label

// SerializeJSON serializes one or more labels into JSON strings and
// appends to the given buffer slice, and then returns the appended
// buffer slice.
func (l Labels) SerializeJSON(buffer []byte) []byte {
	buffer = append(buffer, '{')
	tail := len(l) - 1
	for index := 0; index < len(l); index++ {
		buffer = l[index].SerializeJSON(buffer)
		if index != tail {
			buffer = append(buffer, ", "...)
		}
	}
	return append(buffer, '}')
}

// SerializeStandard serializes one or more labels into standard log
// strings and appends to the given buffer slice, and then returns the
// appended buffer slice.
func (l Labels) SerializeStandard(buffer []byte) []byte {
	return l.SerializeJSON(buffer)
}

// SerializedLabels is a structure that contains data of one or more
// labels that have been serialized using multiple encoding formats to
// avoid repeatedly serializing a set of the same labels.
//
// For details, please refer to the notes section of the Labels structure.
type SerializedLabels struct {
	count int
	jsonBuffer []byte
}

// Count returns the number of labels.
func (l SerializedLabels) Count() int {
	return l.count
}

// SerializeJSON appends a set of serialized label JSON strings to the
// given buffer slice, and then returns the appended buffer slice.
func (l SerializedLabels) SerializeJSON(buffer []byte) []byte {
	return append(buffer, l.jsonBuffer...)
}

// SerializeStandard appends a set of serialized label standard log
// strings to the given buffer slice, and then returns the appended
// buffer slice.
func (l SerializedLabels) SerializeStandard(buffer []byte) []byte {
	return l.SerializeJSON(buffer)
}

// NewSerializedLabels pre-serializes a given set of labels, and then
// returns a SerializedLabels value.
func NewSerializedLabels(labels ...Label) SerializedLabels {
	return SerializedLabels {
		count: len(labels),
		jsonBuffer: Labels(labels).SerializeJSON(make([]byte, 0, 256)),
	}
}
