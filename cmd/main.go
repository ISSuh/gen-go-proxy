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
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ISSuh/simple-gen-proxy/internal/option"
	"github.com/ISSuh/simple-gen-proxy/internal/parser"
)

const (
	proxyFileoutFilePathSuffix = "_proxy"
	sourceFileExtention        = ".go"
	txMiddlewareFileName       = "proxy_middleware_tx"
)

var errGenFailed = errors.New("failed to generate proxy")

type Import struct {
	Alias string
	Path  string
}

func main() {
	args := option.NewArguments()
	if err := args.Validate(); err != nil {
		panic(err)
	}

	// Get target files
	targetDir, err := pathToAbsPath(args.Target)
	if err != nil {
		panic(errors.Join(errGenFailed, err))
	}

	items, err := os.ReadDir(targetDir)
	if err != nil {
		panic(errors.Join(errGenFailed, err))
	}

	fileNames := []string{}
	for _, item := range items {
		if item.IsDir() {
			continue
		}

		fileNames = append(fileNames, item.Name())
	}

	fmt.Printf("Target Dir: %s\n", targetDir)
	fmt.Printf("Files: %v\n", fileNames)

	// Get output path
	outPath, err := outputPath(args.Output, targetDir)
	if err != nil {
		panic(errors.Join(errGenFailed, err))
	}

	// Generate proxy files from target files
	packgeName := args.Package
	for _, fileName := range fileNames {
		filePath := filepath.Join(targetDir, fileName)
		if err != nil {
			panic(errors.Join(errGenFailed, err))
		}

		fileNameWithOutExt := fileNameWithoutExt(fileName)
		outFileName := fileNameWithOutExt + proxyFileoutFilePathSuffix + sourceFileExtention
		outFilePath := filepath.Join(outPath, outFileName)

		// Parse target file
		g := parser.NewGenerator()

		param := parser.ParseParam{
			TargetFile:           filePath,
			TargetFileDir:        targetDir,
			OutFile:              outFileName,
			ProxyPackageName:     args.Package,
			InterfacePackageName: args.InterfacePackage.Name,
			InterfacePackagePath: args.InterfacePackage.Path,
		}

		tmpl, err := g.Parse(param)
		if err != nil {
			panic(errors.Join(errGenFailed, err))
		}

		if len(tmpl.Data.Interfaces) == 0 {
			fmt.Printf("Generate proxy: No interface found in %s\n", fileName)
			continue
		}

		fmt.Printf("Generate proxy: from %v. to %s\n", tmpl.Data.Interfaces.Names(), outFilePath)

		// Generate proxy file
		if err := g.GenerateProxy(outFilePath, tmpl); err != nil {
			panic(errors.Join(errGenFailed, err))
		}

		packgeName = tmpl.Data.PackageName
	}

	// Generate transaction middleware
	if args.UseTxMiddleware {
		outFileName := txMiddlewareFileName + sourceFileExtention
		outFilePath := filepath.Join(outPath, outFileName)

		tmpl := parser.Template{
			Data: &parser.TemplateData{
				PackageName: packgeName,
			},
		}

		g := parser.NewGenerator()
		if err := g.GenerateTxMiddleware(outFilePath, tmpl); err != nil {
			panic(errors.Join(errGenFailed, err))
		}

		fmt.Printf("Generate proxy: generate transaction middleware. To %s\n", outPath)
	}
}

func pathToAbsPath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func outputPath(output, targetDir string) (string, error) {
	outPath := output
	if outPath == "" {
		outPath = targetDir
	}

	outPath, err := pathToAbsPath(outPath)
	if err != nil {
		return "", err
	}

	return outPath, nil
}

func fileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(sourceFileExtention)])
}
