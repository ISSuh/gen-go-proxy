// MIT License

// Copyright (c) 2025 ISSuh

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package service

import (
	"fmt"
)

type Foo interface {
	// use @{annotation name} comment if you want to generate proxy code
	// @proxy
	Logic(needEmitErr bool) (string, error)

	// also support multiple annotation
	// proxy middleware runs in order of annotation
	// @custom1
	// @custom2
	Foo() int
}

type Bar interface {
	// use @{annotation name} comment if you want to generate proxy code
	// @proxy
	Logic(needEmitErr bool) (string, error)

	// also support multiple annotation
	// proxy middleware runs in order of annotation
	// @custom1
	// @custom2
	Foo() int
}

type foo struct{}

func NewFoo() Foo {
	return &foo{}
}

func (f *foo) Logic(needEmitErr bool) (string, error) {
	fmt.Println("[Foo] logic")
	if needEmitErr {
		return "", fmt.Errorf("emit error")
	}
	return "foo logic", nil
}

func (f *foo) Foo() int {
	fmt.Println("[Foo] foo")
	return 1
}

type bar struct{}

func NewBar() Bar {
	return &bar{}
}

func (b *bar) Logic(needEmitErr bool) (string, error) {
	fmt.Println("[Bar] logic")
	if needEmitErr {
		return "", fmt.Errorf("emit error")
	}
	return "bar logic", nil
}

func (b *bar) Foo() int {
	fmt.Println("[Bar] foo")
	return 1
}
