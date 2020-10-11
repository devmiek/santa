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
	"testing"

	"github.com/nobody-night/santa"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	SetName("testing")
	SetLevel(santa.LevelFatal)
	SetSampler(nil)
	SetLabels(santa.NewLabel("name", "testing"))
	AddHooks(santa.NewSimpleHook(func(entry *santa.Entry) error {
		return nil
	}))
	ResetHooks()
	
	err := Duplicate().Close()
	assert.NoError(t, err, "Unexpected close error")
}

func TestStructured(t *testing.T) {
	instance, err := santa.NewStandardBenchmark(false, santa.EncoderStandard)
	assert.NoError(t, err, "Unexpected create error")
	assert.NotNil(t, instance, "Unexpected return value")
	
	err = Set(instance)
	assert.NoError(t, err, "Unexpected set error")

	err = Prints(santa.LevelFatal, "testing", santa.String("name", "testing"))
	assert.NoError(t, err, "Unexpected print error")

	err = Debugs("testing", santa.String("name", "testing"))
	assert.NoError(t, err, "Unexpected print error")

	err = Infos("testing", santa.String("name", "testing"))
	assert.NoError(t, err, "Unexpected print error")

	err = Warnings("testing", santa.String("name", "testing"))
	assert.NoError(t, err, "Unexpected print error")

	err = Errors("testing", santa.String("name", "testing"))
	assert.NoError(t, err, "Unexpected print error")

	err = Fatals("testing", santa.String("name", "testing"))
	assert.NoError(t, err, "Unexpected print error")

	err = Sync()
	assert.NoError(t, err, "Unexpected sync error")

	err = Close()
	assert.NoError(t, err, "Unexpected close error")
}

func TestTemplate(t *testing.T) {
	instance, err := santa.NewStandardBenchmark(false, santa.EncoderStandard)
	assert.NoError(t, err, "Unexpected create error")
	assert.NotNil(t, instance, "Unexpected return value")
	
	err = Set(instance)
	assert.NoError(t, err, "Unexpected set error")

	err = Printf(santa.LevelFatal, "testing %s", "santa")
	assert.NoError(t, err, "Unexpected print error")

	err = Debugf("testing %s", "santa")
	assert.NoError(t, err, "Unexpected print error")

	err = Infof("testing %s", "santa")
	assert.NoError(t, err, "Unexpected print error")

	err = Warningf("testing %s", "santa")
	assert.NoError(t, err, "Unexpected print error")

	err = Errorf("testing %s", "santa")
	assert.NoError(t, err, "Unexpected print error")

	err = Fatalf("testing %s", "santa")
	assert.NoError(t, err, "Unexpected print error")

	err = Sync()
	assert.NoError(t, err, "Unexpected sync error")

	err = Close()
	assert.NoError(t, err, "Unexpected close error")
}
