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

package repository

import (
	"context"
	"database/sql"
	"errors"
)

type TxHepler struct {
	db *sql.DB
}

func NewTxHepler(db *sql.DB) TxHepler {
	return TxHepler{
		db: db,
	}
}

func (t TxHepler) Transaction() func(c context.Context, f func(c context.Context) error) {
	return func(c context.Context, f func(c context.Context) error) {
		requestID, ok := c.Value("requestID").(string)
		if !ok {
			panic(errors.New("requestID not found"))
		}

		txKey := requestID
		tx, ok := c.Value(txKey).(*sql.Tx)
		if tx == nil {
			t.newTransaction(c, f, requestID)
		} else {
			if !ok {
				panic(errors.New("transaction not found"))
			}

			t.subTransaction(c, f, tx)
		}

		tx.Commit()
	}
}

func (t TxHepler) newTransaction(c context.Context, f func(c context.Context) error, txKey string) {
	tx, err := t.db.Begin()
	if err != nil {
		panic(err)
	}

	c = context.WithValue(c, txKey, tx)
	if err := f(c); err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
}

func (t TxHepler) subTransaction(c context.Context, f func(c context.Context) error, tx *sql.Tx) {
	err := f(c)
	if err != nil {
		tx.Rollback()
		return
	}
}
