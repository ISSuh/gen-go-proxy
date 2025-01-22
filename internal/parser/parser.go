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
	"bytes"
	"embed"
	_ "embed"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"os"
)

//go:embed templates/*.tmpl
var templateFiles embed.FS

type TemplateData struct {
	PackageName string
	Imports     []Import
	Interfaces  []Interface
}

type Generator struct {
	fset *token.FileSet
	node *ast.File
}

func NewGenerator() *Generator {
	return &Generator{
		fset: token.NewFileSet(),
	}
}

func (g *Generator) Parse(path, packageName, interfacePakageName, interfacePakagePath string) (*TemplateData, error) {
	node, err := parser.ParseFile(g.fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	g.node = node

	imports, err := ParseImportPackage(node)
	if err != nil {
		return nil, err
	}

	if interfacePakageName != "" && interfacePakagePath != "" {
		imports = append(imports, Import{
			Alias: interfacePakageName,
			Path:  interfacePakagePath,
		})
	}

	iface, err := ParseInterface(node)
	if err != nil {
		return nil, err
	}

	name := packageName
	if packageName == "" {
		name = node.Name.Name
	}

	data := &TemplateData{
		PackageName: name,
		Imports:     imports,
		Interfaces:  iface,
	}
	return data, nil
}

func (g *Generator) Generate(outPath string, data *TemplateData) error {
	tmpl, err := template.ParseFS(templateFiles, "templates/*.tmpl")
	if err != nil {
		panic(err)
	}
	file, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return err
	}

	code, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	if _, err := file.Write(code); err != nil {
		return err
	}

	if err := file.Sync(); err != nil {
		return err
	}

	return nil
}
