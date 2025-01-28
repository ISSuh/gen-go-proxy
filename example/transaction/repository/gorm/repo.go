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

package infra

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

const txKey = "tx"

func FromContext(c context.Context) (*gorm.DB, error) {
	db, ok := c.Value(txKey).(*gorm.DB)
	if !ok || db == nil {
		return nil, errors.New("transaction not found")
	}

	return db, nil
}

// implement example user-defined transaction from transaction interface
type gormTransaction struct {
	db *gorm.DB
}

func NewGORMTransaction(db *gorm.DB) *gormTransaction {
	return &gormTransaction{
		db: db,
	}
}

func (t *gormTransaction) From(c context.Context) error {
	db, ok := c.Value(txKey).(*gorm.DB)
	if !ok {
		return errors.New("transaction not found")
	}

	t.db = db
	return nil
}

func (t *gormTransaction) Begin() error {
	t.db = t.db.Begin()
	return nil
}

func (t *gormTransaction) Commit() error {
	t.db = t.db.Commit()
	return nil
}

func (t *gormTransaction) Rollback() error {
	t.db = t.db.Rollback()
	return nil
}

func (t *gormTransaction) Regist(c context.Context) context.Context {
	return context.WithValue(c, txKey, t.db)
}
