﻿// MIT License

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
	"errors"
	"fmt"
	"go/ast"
	"strings"
	"unicode"
)

const (
	transactionComment = "@transactional"
	proxyComment       = "@proxy"
	annotationToken    = '@'
	minAnnotationLen   = 2

	errorType   = "error"
	contextType = "context.Context"

	userContextParam   = "_userCtx"
	helperContextParam = "_helperCtx"
)

type Annotation struct {
	ProxyTypeName  string
	AnnotationName string
	MethodName     string
}

type Annotations []Annotation

func (a Annotations) Exist(annotation string) bool {
	for _, an := range a {
		if an.AnnotationName == annotation {
			return true
		}
	}
	return false
}

type Param struct {
	Type       string
	Var        string
	HasContext bool
}

func (p Param) Format() string {
	return p.Var + " " + p.Type
}

type Params []Param

func (p Params) HasContext() bool {
	for _, param := range p {
		if param.HasContext {
			return true
		}
	}
	return false
}

func (p Params) Format() string {
	formats := []string{}
	for _, param := range p {
		formats = append(formats, param.Format())
	}
	return strings.Join(formats, ", ")
}

func (p Params) FormatVars(useHelperContext bool) string {
	params := []string{}
	for _, param := range p {
		if useHelperContext && param.Type == contextType {
			params = append(params, helperContextParam)
		} else {
			params = append(params, param.Var)
		}
	}
	return strings.Join(params, ", ")
}

type Result struct {
	ResultType string
	ResultVar  string
}

type Results []Result

func (r Results) FormatType() string {
	if len(r) == 0 {
		return ""
	}
	if len(r) == 1 {
		return r[0].ResultType
	}

	results := []string{}
	for _, result := range r {
		results = append(results, result.ResultType)
	}

	return "(" + strings.Join(results, ", ") + ")"
}

func (r Results) FormatVars() string {
	vars := []string{}
	for _, result := range r {
		vars = append(vars, result.ResultVar)
	}
	return strings.Join(vars, ", ")
}

func (r Results) HasError() bool {
	for _, result := range r {
		if result.ResultType == errorType {
			return true
		}
	}
	return false
}

type Method struct {
	ProxyTypeName               string
	Name                        string
	Annotations                 Annotations
	Params                      string
	ParamNames                  string
	ParamNamesWithHelperContext string
	UserContextParam            string
	HelperContextParam          string
	ResultVars                  string
	ResultTypes                 string
	Results                     Results
	HasResults                  bool
	UseProxy                    bool
	HasError                    bool
	HasContext                  bool
}

type Methods []Method

func (m Methods) AllAnnotations() Annotations {
	annotations := Annotations{}
	for _, method := range m {
		for _, annotation := range method.Annotations {
			if annotations.Exist(annotation.AnnotationName) {
				continue
			}

			annotations = append(annotations, annotation)
		}
	}
	return annotations
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

		annotations := parseAnnotation(method, methodName, proxyTypeName)

		params, err := parseMethodParams(funcType)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("failed to parse method params for %s", methodName), err)
		}

		results, err := parseMethodResults(funcType)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("failed to parse method results for %s", methodName), err)
		}

		m := Method{
			ProxyTypeName: proxyTypeName,
			Name:          methodName,
			Annotations:   annotations,
			UseProxy:      len(annotations) != 0,
			Params:        params.Format(),
			ParamNames:    params.FormatVars(false),
			Results:       results,
			ResultVars:    results.FormatVars(),
			ResultTypes:   results.FormatType(),
			HasResults:    len(results) > 0,
			HasError:      results.HasError(),
			HasContext:    params.HasContext(),
		}

		if m.HasContext {
			m.UserContextParam = userContextParam
			m.HelperContextParam = helperContextParam
			m.ParamNamesWithHelperContext = params.FormatVars(true)
		}

		methods = append(methods, m)
	}

	return methods, nil
}

func isTransactionMethod(method *ast.Field) bool {
	return method.Doc != nil && strings.Contains(method.Doc.Text(), transactionComment)
}

func isProxyMethod(method *ast.Field) bool {
	return method.Doc != nil && strings.Contains(method.Doc.Text(), proxyComment)
}

func parseAnnotation(method *ast.Field, methodName, proxyTypeName string) Annotations {
	if method.Doc == nil {
		return nil
	}

	annotations := Annotations{}
	lines := strings.Split(method.Doc.Text(), "\n")
	for i := range lines {
		index := len(lines) - i - 1
		if !isValidAnnotation(lines[index]) {
			continue
		}

		annotation := strings.ToLower(lines[index][1:])
		if annotations.Exist(annotation) {
			continue
		}

		a := Annotation{
			AnnotationName: strings.ToLower(lines[index][1:]),
			MethodName:     methodName,
			ProxyTypeName:  proxyTypeName,
		}

		annotations = append(annotations, a)
	}
	return annotations
}

func parseMethodParams(funcType *ast.FuncType) (Params, error) {
	params := Params{}
	hasContext := false
	for _, param := range funcType.Params.List {
		for _, name := range param.Names {
			paramName := name.Name
			paramType := exprToString(param.Type)

			if paramType == contextType {
				if hasContext {
					return nil, errors.New("method must have at most one context.Context parameter")
				}

				hasContext = true
				paramName = userContextParam
			}

			p := Param{
				Type:       paramType,
				Var:        paramName,
				HasContext: hasContext,
			}

			params = append(params, p)
		}
	}

	return params, nil
}

func parseMethodResults(funcType *ast.FuncType) (Results, error) {
	results := Results{}
	hasError := false
	if funcType.Results != nil {
		for i, result := range funcType.Results.List {
			resultType := exprToString(result.Type)
			vars := fmt.Sprintf("r%d", i)
			if resultType == errorType {
				if hasError {
					return nil, errors.New("proxy or transactional method must have at most one error result")
				}

				hasError = true
				vars = "err"
			}

			r := Result{
				ResultType: resultType,
				ResultVar:  vars,
			}

			results = append(results, r)
		}
	}

	return results, nil
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
	case *ast.SliceExpr:
		return "[]" + exprToString(t.X)
	case *ast.MapType:
		return "map[" + exprToString(t.Key) + "]" + exprToString(t.Value)
	case *ast.FuncType:
		return exprFuncToString(t)
	case *ast.FuncLit:
		return "func" // Simplified for brevity
	case *ast.Ellipsis:
		return "..." + exprToString(t.Elt)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.StructType:
		return "struct{}"
	case *ast.ChanType:
		return "chan " + exprToString(t.Value)
	case *ast.ParenExpr:
		return "(" + exprToString(t.X) + ")"
	default:
		return ""
	}
}

func exprFuncToString(t *ast.FuncType) string {
	params := []string{}
	for _, param := range t.Params.List {
		paramType := exprToString(param.Type)
		for _, name := range param.Names {
			params = append(params, name.Name+" "+paramType)
		}
		if len(param.Names) == 0 {
			params = append(params, paramType)
		}
	}

	results := []string{}
	if t.Results != nil {
		for _, result := range t.Results.List {
			results = append(results, exprToString(result.Type))
		}
	}

	resultFormat := strings.Join(results, ", ")
	if len(results) > 1 {
		resultFormat = "(" + resultFormat + ")"
	}

	return "func(" + strings.Join(params, ", ") + ") " + resultFormat
}

func isValidAnnotation(s string) bool {
	if len(s) < minAnnotationLen {
		return false
	}

	if s[0] != annotationToken {
		return false
	}

	if strings.Contains(s, " ") {
		return false
	}

	for _, r := range s[1:] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}
