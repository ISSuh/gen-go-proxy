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

package parser

import (
	"go/ast"
)

const (
	proxySuffix = "Proxy"
)

type Interface struct {
	ProxyTypeName     string
	InterfaceName     string
	InterfacePackage  string
	IsDiffrentPackage bool
	Methods           []Method
	types             *ast.InterfaceType
}

type Interfaces []Interface

func (i Interfaces) Names() []string {
	names := []string{}
	for _, iface := range i {
		names = append(names, iface.InterfaceName)
	}
	return names
}

func ParseInterface(node *ast.File, isDiffrentPackage bool) ([]Interface, error) {
	interfaces, err := parseInterfaceType(node, isDiffrentPackage)
	if err != nil {
		return nil, err
	}

	for i, iface := range interfaces {
		interfaces[i].ProxyTypeName = interfaces[i].InterfaceName + proxySuffix
		m, err := parseMethod(interfaces[i].ProxyTypeName, iface.types)
		if err != nil {
			return nil, err
		}

		interfaces[i].Methods = append(interfaces[i].Methods, m...)
	}

	return interfaces, nil
}

func parseInterfaceType(node *ast.File, isDiffrentPackage bool) ([]Interface, error) {
	interfaces := []Interface{}
	ast.Inspect(node, func(n ast.Node) bool {
		spec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		iface, ok := spec.Type.(*ast.InterfaceType)
		if !ok {
			return true
		}

		i := Interface{
			types:             iface,
			InterfaceName:     spec.Name.Name,
			InterfacePackage:  node.Name.Name,
			IsDiffrentPackage: isDiffrentPackage,
		}

		interfaces = append(interfaces, i)
		return true
	})

	return interfaces, nil
}
