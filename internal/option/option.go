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

package option

import (
	"errors"

	"github.com/alexflint/go-arg"
)

type InterfacePackage struct {
	Name string `arg:"-n,--interface-package-name" help:""package name of the target interface source code file"`
	Path string `arg:"-l,--interface-package-path" help:""package name of the target interface source code file"`
}

type Arguments struct {
	InterfacePackage
	Target  string `arg:"-t,--target" help:"target interface source code file"`
	Output  string `arg:"-o,--output" help:"output file path. default is current directory"`
	Suffix  string `arg:"-s,--suffix" help:"suffix for the generated code. default is '_proxy.go'"`
	Pakcage string `arg:"-p,--package" help:"package name of the generated code. default is the same as the target interface source code file"`
}

func NewArguments() Arguments {
	a := Arguments{}
	arg.MustParse(&a)
	return a
}

func (a *Arguments) Validate() error {
	if a.Target == "" {
		return errors.New("target interface source code file is empty")
	}

	if a.Output == "" {
		a.Output = "./"
	}

	return nil
}
