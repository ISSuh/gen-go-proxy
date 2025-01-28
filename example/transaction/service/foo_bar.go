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

package service

import (
	"context"
	"fmt"

	"github.com/ISSuh/gen-go-proxy/example/transaction/dto"
	entity "github.com/ISSuh/gen-go-proxy/example/transaction/entity"
)

type FooBar interface {
	// @transactional
	Create(c context.Context, foo dto.Foo, bar dto.Bar) (int, int, error)

	Find(c context.Context, fooID, barID int) (*entity.Foo, *entity.Bar, error)
}

type fooBarService struct {
	fooService Foo
	barService Bar
}

func NewFooBarService(foo Foo, bar Bar) FooBar {
	return &fooBarService{
		fooService: foo,
		barService: bar,
	}
}

func (s *fooBarService) Create(c context.Context, foo dto.Foo, bar dto.Bar) (int, int, error) {
	fmt.Printf("[FooBar]Create\n")

	fooID, err := s.fooService.Create(c, foo)
	if err != nil {
		return 0, 0, err
	}

	barID, err := s.barService.Create(c, bar)
	if err != nil {
		return 0, 0, err
	}

	return fooID, barID, nil
}

func (s *fooBarService) Find(c context.Context, fooID, barID int) (*entity.Foo, *entity.Bar, error) {
	fmt.Printf("[FooBar]Find\n")

	foo, err := s.fooService.Find(c, fooID)
	if err != nil {
		return nil, nil, err
	}

	bar, err := s.barService.Find(c, barID)
	if err != nil {
		return nil, nil, err
	}

	return foo, bar, nil
}
