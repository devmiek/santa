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

// Exporter is a public interface for exporters.
//
// The exporter uses a specific encoder to encode log entries into
// specific data, and then uses a specific synchronizer to write the
// encoded log entry data to a specific storage device.
//
// Normally, the exporter checks the level of each log entry to determine
// whether it needs to be processed. This means that with exporters,
// different encoders and synchronizers can be set for different log level
// spans. For example, for log entries with levels from DEBUG to DEBUG, use
// a standard encoder and standard synchronizer to output to the standard
// output of the console; for log entries with levels from INFO to WARNING,
// use a JSON encoder and file synchronizer to output to local files, etc.
//
// In addition, the exporter provides temporary buffers and thread safety
// for all encoders and synchronizers.
//
// Please note that every exporter must be closed manually after it is no
// longer used. For details, please refer to the Syncer interface.
type Exporter interface {
	// Export encodes a given log entry into specific data using a specific
	// encoder, then uses a specific synchronizer to write the encoded log
	// entry data to a specific storage device.
	//
	// Finally, any errors encountered are returned.
	Export(entry *Entry) error

	// Sync writes the internal cache data of a specific synchronizer to a
	// specific storage device. If the specific storage device is based on
	// the file system, write the data cached by the file system to the
	// persistent storage device. For details, please refer to the Sync
	// function of the Syncer interface.
	//
	// Finally, any errors encountered are returned.
	Sync() error

	// Close close a specific synchronizer. For details, please participate
	// in the Close function of the Syncer interface.
	//
	// Finally, any errors encountered are returned.
	Close() error
}

// StandardExporter is the structure of the standard exporter instance.
// 
// The standard exporter checks whether the level of each log entry is
// included in the log level span. If it is included, use a specific
// encoder to encode the log entry into specific data, and then use a
// specific synchronizer to write the encoded log entry data to a specific
// storage device.
type StandardExporter struct {
	span LevelSpan
	encoder Encoder
	syncer Syncer
}

// Export encodes a given log entry into specific data using a specific
// encoder, then uses a specific synchronizer to write the encoded log
// entry data to a specific storage device.
//
// Finally, any errors encountered are returned.
func (e *StandardExporter) Export(entry *Entry) error {
	if !e.span.Contains(entry.Level) {
		return nil
	}

	if e.encoder == nil {
		return nil
	}

	pointer := pool.buffer.exporter.New()
	buffer, err := e.encoder.Encode((*pointer)[ : 0], entry)

	if err != nil {
		pool.buffer.exporter.Free(pointer)
		return err
	}

	if buffer == nil {
		pool.buffer.exporter.Free(pointer)
		return nil
	}

	if e.syncer == nil {
		pool.buffer.exporter.Free(pointer)
		return nil
	}

	_, err = e.syncer.Write(buffer)
	pool.buffer.exporter.Free(pointer)

	return err
}

// Sync writes the internal cache data of a specific synchronizer to a
// specific storage device. If the specific storage device is based on
// the file system, write the data cached by the file system to the
// persistent storage device. For details, please refer to the Sync
// function of the Syncer interface.
//
// Finally, any errors encountered are returned.
func (e *StandardExporter) Sync() error {
	return e.syncer.Sync()
}

// Close close a specific synchronizer. For details, please participate
// in the Close function of the Syncer interface.
//
// Finally, any errors encountered are returned.
func (e *StandardExporter) Close() error {
	return e.syncer.Close()
}

// StandardExporterOption is a structure that contains exporter options.
type StandardExporterOption struct {
	// Span represents the log level span. If the level of a log entry is
	// included in the log level span, the log entry will be processed,
	// otherwise it will be discarded. If not provided, the default value
	// is INFO level to FATAL level.
	Span LevelSpan

	// Encoder represents the encoder used to encode log entries. If not
	// provided, the default value is the standard encoder.
	Encoder Encoder
	
	// Syncer represents a synchronizer used to write encoded log entry
	// data to a specific storage device. If not provided, the default
	// value is the standard synchronizer.
	Syncer Syncer
}

// UseSpan uses the given start and end log levels as the value of the
// Span option. For details, please refer to the comment section of the
// Span option. Then return to the option instance itself.
func (o *StandardExporterOption) UseSpan(start, end Level) *StandardExporterOption {
	o.Span = LevelSpan {
		Start: start,
		End: end,
	}

	return o
}

// UseEncoder uses the given encoder as the value of the Encoder option.
// For details, please refer to the comment section of the Encoder option.
// Then return to the option instance itself.
func (o *StandardExporterOption) UseEncoder(encoder Encoder) *StandardExporterOption {
	o.Encoder = encoder
	return o
}

// UseSyncer uses the given syncer as the value of the Syncer option.
// For details, please refer to the comment section of the Syncer option.
// Then return to the option instance itself.
func (o *StandardExporterOption) UseSyncer(syncer Syncer) *StandardExporterOption {
	o.Syncer = syncer
	return o
}

// Build builds and returns a standard exporter instance.
func (o *StandardExporterOption) Build() (*StandardExporter, error) {
	return &StandardExporter {
		span: o.Span,
		encoder: o.Encoder,
		syncer: o.Syncer,
	}, nil
}

// NewStandardExporterOption creates and returns an instance of the
// standard exporter option with default optional values.
func NewStandardExporterOption() *StandardExporterOption {
	// The synchronizer does not need to be closed manually, because
	// the default discard synchronizer writes data to ioutil.Discard.
	//
	// The error is discarded and usually does not occur.
	encoder, _ := NewStandardEncoder()
	syncer, _ := NewDiscardSyncer()

	return &StandardExporterOption {
		Span: LevelSpan {
			Start: LevelDebug,
			End: LevelFatal,
		},
		Encoder: encoder,
		Syncer: syncer,
	}
}

// NewStandardExporter creates and returns an instance of a standard
// exporter using default optional values.
func NewStandardExporter() (*StandardExporter, error) {
	return NewStandardExporterOption().Build()
}
