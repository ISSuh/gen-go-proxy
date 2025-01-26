// source: example/transaction/service/bar.go

package proxy

import (
	"context"

	"github.com/ISSuh/simple-gen-proxy/example/transaction/dto"

	entity "github.com/ISSuh/simple-gen-proxy/example/transaction/entity"

	service "github.com/ISSuh/simple-gen-proxy/example/transaction/service"
)

type BarProxyHelper func(c context.Context, f func(c context.Context) error) error

type BarProxy struct {
	target service.Bar
	tx     BarProxyHelper
}

func NewBarProxy(target service.Bar, tx BarProxyHelper) *BarProxy {
	return &BarProxy{
		target: target,
		tx:     tx,
	}
}

func (p *BarProxy) Create(_userCtx context.Context, dto dto.Bar) (int, error) {
	var (
		r0  int
		err error
	)

	err = p.tx(_userCtx, func(_helperCtx context.Context) error {
		r0, err = p.target.Create(_helperCtx, dto)
		if err != nil {
			return err
		}
		return nil
	})

	return r0, err
}

func (p *BarProxy) Find(_userCtx context.Context, id int) (*entity.Bar, error) {
	return p.target.Find(_userCtx, id)
}
