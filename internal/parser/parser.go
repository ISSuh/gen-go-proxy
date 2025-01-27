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
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"text/template"

	"golang.org/x/tools/imports"
)

const (
	proxyFileoutFilePathSuffix = "_proxy"
	sourceFIleExtention        = ".go"
	proxyTemplatePath          = "templates/target_proxy.go.tmpl"
	txTemplatePath             = "templates/proxy_middleware_tx.go.tmpl"
)

//go:embed templates/target_proxy.go.tmpl
var proxyTemplate embed.FS

//go:embed templates/proxy_middleware_tx.go.tmpl
var txTemplate embed.FS

type TemplateData struct {
	SourceFile  string
	PackageName string
	Imports     []Import
	Interfaces  Interfaces
}

type Template struct {
	FileName string
	FilePath string
	Data     *TemplateData
}

type ParseParam struct {
	TargetFile           string
	TargetFileDir        string
	OutFile              string
	ProxyPackageName     string
	InterfacePackageName string
	InterfacePackagePath string
}

type Generator struct {
}

func NewGenerator() Generator {
	return Generator{}
}

func (g *Generator) Parse(param ParseParam) (Template, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, param.TargetFile, nil, parser.ParseComments)
	if err != nil {
		return Template{}, err
	}

	imports, err := g.paseImport(node, param.InterfacePackageName, param.InterfacePackagePath)
	if err != nil {
		return Template{}, err
	}

	packageName := node.Name.Name
	isDiffrentPackage := false
	if param.ProxyPackageName != "" {
		packageName = param.ProxyPackageName
		isDiffrentPackage = true
	}

	iface, err := ParseInterface(node, isDiffrentPackage)
	if err != nil {
		return Template{}, err
	}

	template := Template{
		FileName: param.OutFile,
		FilePath: param.TargetFileDir,
		Data: &TemplateData{
			SourceFile:  param.TargetFile,
			PackageName: packageName,
			Imports:     imports,
			Interfaces:  iface,
		},
	}
	return template, nil
}

func (g *Generator) paseImport(node *ast.File, interfacePackageName, interfacePackagePath string) ([]Import, error) {
	imports, err := ParseImportPackage(node)
	if err != nil {
		return nil, err
	}

	if interfacePackageName != "" && interfacePackagePath != "" {
		imports = append(imports, Import{
			Alias: interfacePackageName,
			Path:  interfacePackagePath,
		})
	}

	return imports, nil
}

func (g *Generator) GenerateProxy(outFilePath string, tmpl Template) error {
	t, err := template.ParseFS(proxyTemplate, proxyTemplatePath)
	if err != nil {
		return err
	}

	if err := g.generateFile(outFilePath, tmpl, t); err != nil {
		return err
	}

	return nil
}

func (g *Generator) GenerateTxMiddleware(outFilePath string, tmpl Template) error {
	t, err := template.ParseFS(txTemplate, txTemplatePath)
	if err != nil {
		return err
	}

	if err := g.generateFile(outFilePath, tmpl, t); err != nil {
		return err
	}

	return nil
}

func (g *Generator) generateFile(outFilePath string, tmpl Template, t *template.Template) error {
	file, err := os.Create(outFilePath)
	if err != nil {
		err = errors.Join(fmt.Errorf("failed to create file(%s)", outFilePath), err)
		return err
	}
	defer file.Close()

	// generate code from template
	var buf bytes.Buffer
	err = t.Execute(&buf, tmpl.Data)
	if err != nil {
		return err
	}

	// format code
	format, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	// remove unused imports
	code, err := imports.Process("", format, nil)
	if err != nil {
		return err
	}

	// write code to file
	if _, err := file.Write(code); err != nil {
		return err
	}

	if err := file.Sync(); err != nil {
		return err
	}

	return nil
}
