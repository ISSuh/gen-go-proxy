package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"
)

type Import struct {
	Alias string
	Path  string
}

type Method struct {
	ProxyTypeName string
	Name          string
	Params        string
	ParamNames    string
	Results       string
	ResultVars    string
	HasResults    bool
	IsTransaction bool
}

type Interface struct {
	InterfaceName string
	ProxyTypeName string
	Methods       []Method
}

type TemplateData struct {
	PackageName string
	Imports     []Import
	Interfaces  []Interface
}

func main() {
	// 소스 파일을 파싱하여 인터페이스 정보를 추출합니다.
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "/Users/issuh/workspace/git/issuh/simple-gen-proxy/cmd/example/service/service.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	imports := []Import{}
	interfaces := []Interface{}

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

	// AST를 탐색하여 인터페이스를 찾습니다.
	ast.Inspect(node, func(n ast.Node) bool {
		if ts, ok := n.(*ast.TypeSpec); ok {
			if iface, ok := ts.Type.(*ast.InterfaceType); ok {
				ifaceName := ts.Name.Name
				methods := []Method{}
				for _, method := range iface.Methods.List {
					if len(method.Names) == 0 {
						continue
					}
					methodName := method.Names[0].Name
					funcType := method.Type.(*ast.FuncType)

					isTransaction := method.Doc != nil && strings.Contains(method.Doc.Text(), "transaction")

					params := []string{}
					paramNames := []string{}
					for _, param := range funcType.Params.List {
						for _, name := range param.Names {
							paramName := name.Name
							paramType := exprToString(param.Type)
							params = append(params, paramName+" "+paramType)
							paramNames = append(paramNames, paramName)
						}
					}

					results := []string{}
					resultVars := []string{}
					if funcType.Results != nil {
						for i, result := range funcType.Results.List {
							resultType := exprToString(result.Type)
							results = append(results, resultType)
							resultVars = append(resultVars, fmt.Sprintf("r%d", i))
						}
					}

					methods = append(methods, Method{
						ProxyTypeName: ifaceName + "Proxy",
						Name:          methodName,
						Params:        strings.Join(params, ", "),
						ParamNames:    strings.Join(paramNames, ", "),
						Results:       formatResults(results),
						ResultVars:    strings.Join(resultVars, ", "),
						HasResults:    len(results) > 0,
						IsTransaction: isTransaction,
					})
				}

				interfaces = append(interfaces, Interface{
					InterfaceName: ifaceName,
					ProxyTypeName: ifaceName + "Proxy",
					Methods:       methods,
				})
			}
		}
		return true
	})

	data := TemplateData{
		PackageName: node.Name.Name,
		Imports:     imports,
		Interfaces:  interfaces,
	}

	tmpl, err := template.ParseFiles("proxy_template.go.tmpl")
	if err != nil {
		panic(err)
	}

	file, err := os.Create("proxy_generated.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}

	code, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	file.Write(code)
}

func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	case *ast.MapType:
		return "map[" + exprToString(t.Key) + "]" + exprToString(t.Value)
	case *ast.FuncType:
		return "func" // Simplified for brevity
	default:
		return ""
	}
}

func formatResults(results []string) string {
	if len(results) == 0 {
		return ""
	}
	if len(results) == 1 {
		return results[0]
	}
	return "(" + strings.Join(results, ", ") + ")"
}
