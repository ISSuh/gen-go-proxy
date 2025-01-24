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

	entity "github.com/ISSuh/simple-gen-proxy/example/entity"
	"github.com/ISSuh/simple-gen-proxy/example/repository"
)

type Bar interface {
	// @transactional
	CreateC(c context.Context, id int) (*entity.Foo, error)

	// @transactional
	CreateD(c context.Context, id int) (*entity.Bar, error)

	CBarz(c context.Context, id int) (*entity.Foo, *entity.Bar, error)

	// @transactional
	BFoos(c context.Context, a *entity.Foo, b *entity.Bar) error

	QFoosBarz(c context.Context)
}

type BarService struct {
	repo repository.Bar
}

func NewBarService(repo repository.Bar) Bar {
	return &BarService{
		repo: repo,
	}
}

func (s *BarService) CreateC(c context.Context, id int) (*entity.Foo, error) {
	return &entity.Foo{}, nil
}

func (s *BarService) CreateD(c context.Context, id int) (*entity.Bar, error) {
	return &entity.Bar{}, nil
}

func (s *BarService) CBarz(c context.Context, id int) (*entity.Foo, *entity.Bar, error) {
	return &entity.Foo{}, &entity.Bar{}, nil
}

func (s *BarService) BFoos(c context.Context, a *entity.Foo, b *entity.Bar) error {
	return nil
}

func (s *BarService) QFoosBarz(c context.Context) {}
