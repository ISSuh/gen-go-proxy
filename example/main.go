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

package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/ISSuh/simple-gen-proxy/example/entity"
	"github.com/ISSuh/simple-gen-proxy/example/repository"
	"github.com/ISSuh/simple-gen-proxy/example/service"
	"github.com/ISSuh/simple-gen-proxy/example/service/proxy"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CustomLogger struct {
	logger.Interface
}

func (c CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	log.Printf("[INFO] "+msg, data...)
}

func (c CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	log.Printf("[WARN] "+msg, data...)
}

func (c CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	log.Printf("[ERROR] "+msg, data...)
}

func (c CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	log.Printf("[TRACE] %s [%.2fms] [rows:%v] %s", sql, float64(elapsed.Nanoseconds())/1e6, rows, err)
}

func createTestDB(db *gorm.DB) error {
	err := db.AutoMigrate(&entity.Foo{}, &entity.Bar{})
	if err != nil {
		return err
	}
	return nil
}

func openDB() (*gorm.DB, error) {
	dsn := "root:1@tcp(127.0.0.1:33551)/TESTDB"

	config := &gorm.Config{
		Logger: CustomLogger{},
		// Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return nil, err
	}

	if err := createTestDB(db); err != nil {
		return nil, err
	}
	return db, nil
}

type Server struct {
	foo service.Foo
	bar service.Bar
}

func newServer() *Server {
	return &Server{}
}

func (s *Server) init() error {
	db, err := openDB()
	if err != nil {
		return err
	}

	f := func() *sql.DB {
		rawDB, err := db.DB()
		if err != nil {
			return nil
		}

		return rawDB
	}

	fooRepo := repository.NewFooRepository(db)
	fooService := service.NewFooService(fooRepo)
	s.foo = proxy.NewFooProxy(fooService, f)

	barRepo := repository.NewBarRepository(db)
	barService := service.NewBarService(barRepo)
	s.bar = proxy.NewBarProxy(barService, f)

	s.foo.Barz(context.Background(), 1)

	s.foo.CreateB(context.Background(), 100)
	return nil
}

func (s *Server) Run() {
	http.HandleFunc("GET /foo", func(w http.ResponseWriter, r *http.Request) {
		s.foo.CreateB(nil, 1)
	})

	http.ListenAndServe(":8080", nil)
}

func main() {
	server := newServer()
	if err := server.init(); err != nil {
		panic(err)
	}

	server.Run()
}
