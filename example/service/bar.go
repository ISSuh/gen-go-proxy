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

type Bar interface {
	// @transactional
	Create(c context.Context, dto dto.Bar) (int, error)

	Find(c context.Context, id int) (*entity.Bar, error)
}

type barService struct {
	repo repository.Bar
}

func NewBarService(repo repository.Bar) Bar {
	return &barService{
		repo: repo,
	}
}

func (s *barService) Create(c context.Context, dto dto.Bar) (int, error) {
	fmt.Printf("[Bar]CreateA\n")
	id, err := s.repo.Create(c, dto.Value)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *barService) Find(c context.Context, id int) (*entity.Bar, error) {
	fmt.Printf("[Bar]Find\n")
	bar, err := s.repo.Find(id)
	if err != nil {
		return nil, err
	}
	return bar, nil
}
