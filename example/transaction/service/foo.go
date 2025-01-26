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

	"github.com/ISSuh/simple-gen-proxy/example/transaction/dto"
	entity "github.com/ISSuh/simple-gen-proxy/example/transaction/entity"
	"github.com/ISSuh/simple-gen-proxy/example/transaction/repository"
)

type Foo interface {
	// @transactional
	Create(c context.Context, dto dto.Foo) (int, error)

	Find(c context.Context, id int) (*entity.Foo, error)

	// @transactional
	FooBara(c context.Context, dto dto.Foo) error
}

type FooService struct {
	repo repository.Foo
}

func NewFooService(repo repository.Foo) Foo {
	return &FooService{
		repo: repo,
	}
}

func (s *FooService) Create(c context.Context, dto dto.Foo) (int, error) {
	fmt.Printf("[Foo]CreateA\n")
	id, err := s.repo.Create(c, dto.Value)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *FooService) Find(c context.Context, id int) (*entity.Foo, error) {
	fmt.Printf("[Foo]Find\n")
	foo, err := s.repo.Find(id)
	if err != nil {
		return nil, err
	}
	return foo, nil
}

func (s *FooService) FooBara(c context.Context, dto dto.Foo) error {
	fmt.Printf("[Foo]FooBara\n")
	return nil
}
