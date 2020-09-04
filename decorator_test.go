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

func TestDecorator(t *testing.T) {
	logger, err := New()
	assert.NoError(t, err, "Unexpected build error")

	decorator := logger.Decorator()
	assert.NotNil(t, decorator, "Unexpected return value")
	assert.Equal(t, LevelDebug, decorator.level, "Unexpected logger level")

	decorator.SetLevel(LevelFatal)
	assert.Equal(t, LevelFatal, decorator.level, "Unexpected logger level")

	decorator.SetName("testing")
	assert.Equal(t, "testing", decorator.name, "Unexpected logger name")

	decorator.SetHooks(&SimpleHook { })
	assert.Len(t, decorator.hooks, 1, "Unexpected hook chain length")

	decorator.SetExporters(&StandardExporter { })
	assert.Len(t, decorator.exporters, 1, "Unexpected exporter chain length")

	decorator.SetLabels(NewLabel("testing", "testing"))
	assert.Equal(t, decorator.labels.Count(), 1, "Unexpected exporter chain length")

	decorator.UseHooks(&SimpleHook { })
	assert.Len(t, decorator.hooks, 2, "Unexpected hook chain length")

	decorator.UseExporters(&StandardExporter { })
	assert.Len(t, decorator.exporters, 2, "Unexpected exporter chain length")

	decorator.SetSampler(&TextSampler { })
	assert.NotNil(t, decorator.sampler, "Unexpected sampler value")
}

func TestStandardDecorator(t *testing.T) {
	logger, err := NewStandard()
	assert.NoError(t, err, "Unexpected build error")

	decorator := logger.Decorator()
	assert.NotNil(t, decorator, "Unexpected return value")
	assert.Equal(t, LevelDebug, decorator.level, "Unexpected logger level")
}

func TestTemplateDecorator(t *testing.T) {
	logger, err := NewTemplate()
	assert.NoError(t, err, "Unexpected build error")

	decorator := logger.Decorator()
	assert.NotNil(t, decorator, "Unexpected return value")
	assert.Equal(t, LevelDebug, decorator.level, "Unexpected logger level")
}

func TestStructDecorator(t *testing.T) {
	logger, err := NewStruct()
	assert.NoError(t, err, "Unexpected build error")

	decorator := logger.Decorator()
	assert.NotNil(t, decorator, "Unexpected return value")
	assert.Equal(t, LevelDebug, decorator.level, "Unexpected logger level")
}
