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
	"database/sql"
	"errors"

	"github.com/ISSuh/gen-go-proxy/example/transaction/entity"
)

type FooSQLRepository struct {
	db *sql.DB
}

func NewFooSQLRepository(db *sql.DB) *FooSQLRepository {
	return &FooSQLRepository{
		db: db,
	}
}

func (r *FooSQLRepository) Create(c context.Context, value int) (int, error) {
	tx, err := FromContext(c)
	if err != nil {
		return 0, err
	}

	if value < 0 {
		return 0, errors.New("value must be greater than 0")
	}

	result, err := tx.ExecContext(c, "INSERT INTO foo (value) VALUES (?)", value)
	if err != nil {
		return 0, err
	}

	newID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(newID), nil
}

func (r *FooSQLRepository) Find(id int) (*entity.Foo, error) {
	fooID := 0
	fooValue := 0
	err := r.db.QueryRow("SELECT * FROM foo WHERE id = ?", id).Scan(&fooID, &fooValue)
	if err != nil {
		return nil, err
	}

	return &entity.Foo{
		ID:    fooID,
		Value: fooValue,
	}, nil
}
