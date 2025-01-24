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

	"github.com/ISSuh/simple-gen-proxy/example/dto"
	entity "github.com/ISSuh/simple-gen-proxy/example/entity"
	"github.com/ISSuh/simple-gen-proxy/example/repository"
)

type Foo interface {
	// @transactional
	CreateA(c context.Context, id int, dto dto.ADTO) (*entity.Foo, error)

	// @transactional
	CreateB(c context.Context, id int) (*entity.Bar, error)

	Barz(c context.Context, id int) (*entity.Foo, *entity.Bar, error)

	// @transactional
	Foos(c context.Context, a *entity.Foo, b *entity.Bar) error

	FoosBarz(c context.Context)
}

type FooService struct {
	repo repository.Foo
}

func NewFooService(repo repository.Foo) Foo {
	return &FooService{
		repo: repo,
	}
}

func (s *FooService) CreateA(c context.Context, id int, dto dto.ADTO) (*entity.Foo, error) {
	fmt.Printf("CreateA\n")
	return &entity.Foo{}, nil
}

func (s *FooService) CreateB(c context.Context, id int) (*entity.Bar, error) {
	fmt.Printf("CreateB\n")
	s.repo.Create(int64(id))
	return &entity.Bar{}, nil
}

func (s *FooService) Barz(c context.Context, id int) (*entity.Foo, *entity.Bar, error) {
	fmt.Printf("Barz\n")
	s.repo.Find(id)
	return &entity.Foo{}, &entity.Bar{}, nil
}

func (s *FooService) Foos(c context.Context, a *entity.Foo, b *entity.Bar) error {
	fmt.Printf("Foos\n")
	return nil
}

func (s *FooService) FoosBarz(c context.Context) {
	fmt.Printf("FoosBarz\n")
}
