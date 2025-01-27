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
	"context"
	"fmt"

	"github.com/ISSuh/simple-gen-proxy/example/proxy/service"
)

// implement middleware
func Wrapped(next func(c context.Context) error) func(context.Context) error {
	return func(c context.Context) error {
		fmt.Println("[Wrapped] before")

		// run next middleware or target logic
		err := next(c)
		if err != nil {
			fmt.Printf("[Wrapped] err occurred. err : %s\n", err)
		}

		fmt.Println("[Wrapped] after")
		return err
	}
}

func Before(next func(c context.Context) error) func(context.Context) error {
	return func(c context.Context) error {
		fmt.Println("[Before] before")

		// run next middleware or target logic
		return next(c)
	}
}

func After(next func(c context.Context) error) func(context.Context) error {
	return func(c context.Context) error {
		// run next middleware or target logic
		err := next(c)
		if err != nil {
			fmt.Printf("[After] err occurred. err : %s\n", err)
		}

		fmt.Println("[After] after")
		return err
	}
}

func main() {
	target := service.NewFoo()

	// middleware by annotation
	// key: annotation name
	// value: middleware list
	m := service.FooProxyMiddlewareByAnnotation{
		"proxy":   {Wrapped, Before, After},
		"custom1": {Wrapped},
		"custom2": {Before, After},
	}

	proxy := service.NewFooProxy(target, m)

	if val, err := proxy.Logic(false); err != nil {
		fmt.Println("err: ", err)
	} else {
		fmt.Println("val: ", val)
	}

	fmt.Println()

	if val, err := proxy.Logic(true); err != nil {
		fmt.Println("err: ", err)
	} else {
		fmt.Println("val: ", val)
	}

	fmt.Println()

	value := proxy.Foo()
	fmt.Println("value: ", value)
}
