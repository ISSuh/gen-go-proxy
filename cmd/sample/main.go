package main

import (
	"fmt"
	"go/ast"
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
	Name       string
	Params     string
	ParamNames string
	Results    string
	ResultVars string
	HasResults bool
}

type TemplateData struct {
	PackageName   string
	Imports       []Import
	InterfaceName string
	ProxyTypeName string
	Methods       []Method
}

func main() {
	// 소스 파일을 파싱하여 인터페이스 정보를 추출합니다.
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "source.go", nil, parser.AllErrors)
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

	methods := []Method{}
	for _, method := range ifaceType.Methods.List {
		if len(method.Names) == 0 {
			continue
		}
		methodName := method.Names[0].Name
		funcType := method.Type.(*ast.FuncType)

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
			Name:       methodName,
			Params:     strings.Join(params, ", "),
			ParamNames: strings.Join(paramNames, ", "),
			Results:    formatResults(results),
			ResultVars: strings.Join(resultVars, ", "),
			HasResults: len(results) > 0,
		})
	}

	data := TemplateData{
		PackageName:   node.Name.Name,
		Imports:       imports,
		InterfaceName: ifaceName,
		ProxyTypeName: ifaceName + "Proxy",
		Methods:       methods,
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

	err = tmpl.Execute(file, data)
	if err != nil {
		panic(err)
	}
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
