package service

import (
	"context"

	"github.com/ISSuh/simple-gen-proxy/cmd/example/entity"
)

type Foo interface {
	CreateA(c context.Context, id int) (*entity.A, error)
	CreateB(c context.Context, id int) (*entity.B, error)
	Barz(c context.Context, id int) (*entity.A, *entity.B, error)
	Foos(c context.Context, a *entity.A, b *entity.B) error
}

type FooService struct{}

func NewFooService() Foo {
	return &FooService{}
}

func (s *FooService) CreateA(c context.Context, id int) (*entity.A, error) {
	return &entity.A{}, nil
}

func (s *FooService) CreateB(c context.Context, id int) (*entity.B, error) {
	return &entity.B{}, nil
}

func (s *FooService) Barz(c context.Context, id int) (*entity.A, *entity.B, error) {
	return &entity.A{}, &entity.B{}, nil
}

func (s *FooService) Foos(c context.Context, a *entity.A, b *entity.B) error {
	return nil
}
