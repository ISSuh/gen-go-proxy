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
	"html/template"
	"os"
	"path/filepath"

	"github.com/ISSuh/simple-gen-proxy/internal/option"
	"golang.org/x/tools/imports"
)

const (
	proxyFileoutFilePathSuffix = "_proxy"
	sourceFIleExtention        = ".go"
	embedTemplatePath          = "templates/*.tmpl"
)

//go:embed templates/*.tmpl
var templateFiles embed.FS

type TemplateData struct {
	SourceFile  string
	PackageName string
	Imports     []Import
	Interfaces  []Interface
}

type Proxy struct {
	FileName     string
	FilePath     string
	TemplateData *TemplateData
}

type Generator struct {
	fset *token.FileSet
}

func NewGenerator() *Generator {
	return &Generator{
		fset: token.NewFileSet(),
	}
}

func (g *Generator) Parse(args option.Arguments) (Proxy, error) {
	fmt.Printf("Parsing source file: %s\n", args.Target)

	node, err := parser.ParseFile(g.fset, args.Target, nil, parser.ParseComments)
	if err != nil {
		return Proxy{}, err
	}

	imports, err := g.paseImport(node, args.InterfacePackage.Name, args.InterfacePackage.Path)
	if err != nil {
		return Proxy{}, err
	}

	packageName := node.Name.Name
	isDiffrentPackage := false
	if args.Pakcage != "" {
		packageName = args.Pakcage
		isDiffrentPackage = true
	}

	iface, err := ParseInterface(node, isDiffrentPackage)
	if err != nil {
		return Proxy{}, err
	}

	filePath, fileName, err := separateFileNameFromPath(args.Target)
	if err != nil {
		return Proxy{}, err
	}

	proxy := Proxy{
		FileName: fileName + proxyFileoutFilePathSuffix + sourceFIleExtention,
		FilePath: filePath,
		TemplateData: &TemplateData{
			SourceFile:  args.Target,
			PackageName: packageName,
			Imports:     imports,
			Interfaces:  iface,
		},
	}
	return proxy, nil
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

func (g *Generator) Generate(outPath string, data Proxy) error {
	tmpl, err := template.ParseFS(templateFiles, embedTemplatePath)
	if err != nil {
		panic(err)
	}

	outDirPath := outPath
	if outDirPath == "" {
		outDirPath = data.FilePath
	}

	outDirPath, err = filepath.Abs(outDirPath)
	if err != nil {
		return err
	}

	// create file
	outFilePath := filepath.Join(outDirPath, data.FileName)
	fmt.Printf("Generating proxy file: %s\n", outFilePath)

	file, err := os.Create(outFilePath)
	if err != nil {
		err = errors.Join(fmt.Errorf("failed to create file(%s)", outFilePath), err)
		return err
	}
	defer file.Close()

	// generate code from template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data.TemplateData)
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

func separateFileNameFromPath(path string) (string, string, error) {
	if filepath.Ext(path) != sourceFIleExtention {
		return "", "", errors.New("file is not go source file")
	}

	filePath := filepath.Dir(path)
	if filePath == "" {
		return "", "", errors.New("file path is empty")
	}

	fileName := filepath.Base(path)
	if fileName == "" {
		return "", "", errors.New("file name is empty")
	}

	name := fileName[:len(fileName)-len(sourceFIleExtention)]
	return filePath, name, nil
}
