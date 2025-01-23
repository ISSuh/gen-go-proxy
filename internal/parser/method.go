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
	"fmt"
	"go/ast"
	"strings"
)

const (
	transactionComment = "@transactional"
	errorType          = "error"
)

type Result struct {
	ResultType string
	ResultVar  string
}

type Method struct {
	ProxyTypeName string
	Name          string
	Params        string
	ParamNames    string
	Results       string
	ResultVars    string
	Result        []Result
	HasResults    bool
	IsTransaction bool
	HasError      bool
}

func parseMethod(proxyTypeName string, iface *ast.InterfaceType) ([]Method, error) {
	methods := []Method{}
	for _, method := range iface.Methods.List {
		if len(method.Names) == 0 {
			continue
		}

		methodName := method.Names[0].Name
		funcType, ok := method.Type.(*ast.FuncType)
		if !ok {
			return nil, fmt.Errorf("method %s is not a function", methodName)
		}

		isTransaction := isTransactionMethod(method)
		params, paramNames := parseMethodParams(funcType)
		results, resultVars, resultsSli, hasError := parseMethodResults(funcType)

		m := Method{
			ProxyTypeName: proxyTypeName,
			Name:          methodName,
			IsTransaction: isTransaction,
			Params:        params,
			ParamNames:    paramNames,
			Results:       results,
			ResultVars:    resultVars,
			Result:        resultsSli,
			HasResults:    len(results) > 0,
			HasError:      hasError,
		}

		methods = append(methods, m)
	}

	return methods, nil
}

func isTransactionMethod(method *ast.Field) bool {
	return method.Doc != nil && strings.Contains(method.Doc.Text(), transactionComment)
}

func parseMethodParams(funcType *ast.FuncType) (string, string) {
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

	paramStr := strings.Join(params, ", ")
	paramNameStr := strings.Join(paramNames, ", ")
	return paramStr, paramNameStr
}

func parseMethodResults(funcType *ast.FuncType) (string, string, []Result, bool) {
	resultSli := []Result{}
	results := []string{}
	resultVars := []string{}
	hasError := false
	if funcType.Results != nil {
		for i, result := range funcType.Results.List {
			resultType := exprToString(result.Type)
			results = append(results, resultType)

			vars := fmt.Sprintf("r%d", i)
			if resultType == errorType {
				hasError = true
				// vars = fmt.Sprintf("err%d", i)
				vars = "err"
			}

			resultSli = append(resultSli, Result{ResultType: resultType, ResultVar: vars})
			resultVars = append(resultVars, vars)
		}
	}

	resultStr := formatResults(results)
	resultVarStr := strings.Join(resultVars, ", ")
	return resultStr, resultVarStr, resultSli, hasError
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
