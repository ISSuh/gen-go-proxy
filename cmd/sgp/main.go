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

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/ISSuh/simple-gen-proxy/internal/option"
)

type Import struct {
	Alias string
	Path  string
}

func main() {
	arg := option.NewArguments()
	if err := arg.Validate(); err != nil {
		panic(err)
	}

	fmt.Printf("Target: %s\n", arg.Target)
	fmt.Printf("Output: %s\n", arg.Output)
	fmt.Printf("Suffix: %s\n", arg.Suffix)

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, arg.Target, nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	var ifaceType *ast.InterfaceType
	var ifaceName string
	imports := []Import{}

	// AST를 탐색하여 인터페이스를 찾습니다.
	ast.Inspect(node, func(n ast.Node) bool {
		if ts, ok := n.(*ast.TypeSpec); ok {
			if iface, ok := ts.Type.(*ast.InterfaceType); ok {
				ifaceType = iface
				ifaceName = ts.Name.Name
				return false
			}
		}
		return true
	})

	if ifaceType == nil {
		panic("No interface found in source.go")
	}

	// AST를 탐색하여 import 경로를 찾습니다.
	ast.Inspect(node, func(n ast.Node) bool {
		if imp, ok := n.(*ast.ImportSpec); ok {
			alias := ""
			if imp.Name != nil {
				alias = imp.Name.Name
			}
			imports = append(imports, Import{
				Alias: alias,
				Path:  strings.Trim(imp.Path.Value, "\""),
			})
		}
		return true
	})

	fmt.Printf("Interface Name: %s\n", ifaceName)
	fmt.Printf("Imports: %v\n", imports)
}
