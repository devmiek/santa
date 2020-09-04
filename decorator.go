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

// brush provides an API for changing the copy of the packaged logger
// instance to avoid ambiguity.
type brush struct {
	logger *Logger
}

// SetLevel sets the given log level to the level of the logger.
func (b brush) SetLevel(level Level) {
	b.logger.level = level
}

// SetName sets the given name to the name of the logger.
func (b brush) SetName(name string) {
	b.logger.name = name
}

// SetSampler sets the given sampler instance as the sampler of the
// logger. If the value of a given sampler instance is nil, no sampler
// is useb.logger.
func (b brush) SetSampler(sampler Sampler) {
	b.logger.sampler = sampler
}

// SetHooks sets the given one or more Hook instances as the Hook chain
// of the logger. If no Hook instance is given, no Hook is useb.logger.
func (b brush) SetHooks(hooks ...Hook) {
	b.logger.hooks = hooks
}

// UseHooks appends one or more given Hook instances to the Hook chain
// of the logger.
func (b brush) UseHooks(hooks ...Hook) {
	b.logger.hooks = append(b.logger.hooks, hooks...)
}

// SetExporters sets one or more exporters as the exporter chain of the
// logger. If any exporter instance is provided, no exporter will be
// used by default.
func (b brush) SetExporters(exporters ...Exporter) {
	b.logger.exporters = exporters
}

// UseExporters appends one or more exporters to the exporter chain of
// logger.
func (b brush) UseExporters(exporters ...Exporter) {
	b.logger.exporters = append(b.logger.exporters, exporters...)	
}

// SetLabels sets one or more labels as the relevant labels of the
// logger.
func (b brush) SetLabels(labels ...Label) {
	b.logger.labels = NewSerializedLabels(labels...)
}

// Decorator is the structure of the decorator instance.
//
// The decorator is a wrapper for a copy of the logger, so the decorator
// instance can be considered as an instance of the logger. Normally, the
// logger is read-only, which means that after one or more logger
// instances are successfully constructed, the internal changes cannot be
// made. This is to ensure the integrity of the logger instance.
//
// Unlike the logger, the decorator allows internal changes to its
// instance during the life cycle, including but not limited to changing
// the logger's level, name, label, sampler, and exporter. It is worth
// noting that the decorator is a wrapper for a copy of the logger
// instance, which means that changes to the decorator instance will not
// affect the original logger instance.
//
// The decorator is like the shadow of the logger, and changes to the
// shadow will not affect the original logger. For example, in a
// production environment, you may need to use different logger levels
// based on variables such as settings, services, and user IDs.
// Therefore, when these variables change, you can create a shadow for
// the logger and change the shadow of the logger.
//
// Please note that the API provided by the decorator is not
// thread-safe. Do not share a decorator instance between multiple
// thread contexts. The end of the life cycle of each decorator must
// be earlier than the original logger instance, otherwise the behavior
// is undefineb.logger.
type Decorator struct {
	brush
	Logger
}

// Free returns the decorator. After the refund, the decorator is not
// allowed to be used again, otherwise the behavior is undefined.
func (d *Decorator) Free() {
	pool.decorator.base.Free(d)
}

// StandardDecorator is the structure of an instance of a standard
// decorator.
//
// The standard decorator is used to wrap a copy of the instance of the
// standard logger. For details, please refer to the comment section of
// the Decorator structure.
type StandardDecorator struct {
	brush
	StandardLogger
}

// Free returns the decorator. After the refund, the decorator is not
// allowed to be used again, otherwise the behavior is undefined.
func (d *StandardDecorator) Free() {
	pool.decorator.standard.Free(d)
}

// TemplateDecorator is the structure of an instance of a template
// decorator.
//
// The template decorator is used to wrap a copy of the instance of the
// template logger. For details, please refer to the comment section of
// the Decorator structure.
type TemplateDecorator struct {
	brush
	TemplateLogger
}

// Free returns the decorator. After the refund, the decorator is not
// allowed to be used again, otherwise the behavior is undefined.
func (d *TemplateDecorator) Free() {
	pool.decorator.template.Free(d)
}

// StructDecorator is the structure of an instance of a struct
// decorator.
//
// The struct decorator is used to wrap a copy of the instance of the
// struct logger. For details, please refer to the comment section of
// the Decorator structure.
type StructDecorator struct {
	brush
	StructLogger
}

// Free returns the decorator. After the refund, the decorator is not
// allowed to be used again, otherwise the behavior is undefined.
func (d *StructDecorator) Free() {
	pool.decorator.structure.Free(d)
}
