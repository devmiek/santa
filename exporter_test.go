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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStandardExporterExport(t *testing.T) {
	exporter, err := NewStandardExporter()
	assert.NoError(t, err, "Unexpected creation error")

	err = exporter.Export(entry)
	assert.NoError(t, err, "Unexpected export error")

	err = exporter.Sync()
	assert.NoError(t, err, "Unexpected sync error")

	err = exporter.Close()
	assert.NoError(t, err, "Unexpected close error")
}

func TestStandardExporterOption(t *testing.T) {
	encoder, _ := NewStandardEncoder()
	syncer, _ := NewDiscardSyncer()

	option := NewStandardExporterOption()

	span := LevelSpan {
		Start: LevelInfo,
		End: LevelWarning,
	}

	option.UseEncoder(encoder)
	option.UseSyncer(syncer)
	option.UseSpan(span.Start, span.End)

	assert.Equal(t, encoder, option.Encoder,
		"Unexpected option value")

	assert.Equal(t, syncer, option.Syncer,
		"Unexpected option value")

	assert.Equal(t, span, option.Span,
		"Unexpected option value")

	exporter, err := option.Build()
	assert.NoError(t, err, "Unexpected build error")

	assert.Equal(t, encoder, exporter.encoder,
		"Unexpected instance error")
	
	assert.Equal(t, syncer, exporter.syncer,
		"Unexpected instance error")

	assert.Equal(t, span, exporter.span,
		"Unexpected instance error")
}
