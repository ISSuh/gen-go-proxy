# transaction example

This example guides you on how to generate and use proxy code and transaction middleware with simple-gen-proxy.

## Usage

### implement interface on service layer

```go
package service

type Foo interface {
  // use @transactional comment if you want to generate proxy code
  // @transactional
  Create(c context.Context, dto dto.Foo) (int, error)

  Find(c context.Context, id int) (*entity.Foo, error)
}

type Bar interface {
  // @transactional
  Create(c context.Context, dto dto.Bar) (int, error)

  Find(c context.Context, id int) (*entity.Bar, error)
}

type FooBar interface {
  // @transactional
  Create(c context.Context, foo dto.Foo, bar dto.Bar) (int, int, error)

  Find(c context.Context, fooID, barID int) (*entity.Foo, *entity.Bar, error)
}
```

### implement user trasaction on repository layer

```go
package repository

import (
  "context"
  "database/sql"
  "errors"
)

const txKey = "tx"

type sqlTransaction struct {
  db *sql.DB
  tx *sql.Tx
}

func NewSQLTransaction(db *sql.DB) *sqlTransaction {
  return &sqlTransaction{
    db: db,
  }
}

func (t *sqlTransaction) From(c context.Context) error {
  tx, ok := c.Value(txKey).(*sql.Tx)
  if !ok {
    return errors.New("transaction not found")
  }

  t.tx = tx
  return nil
}

func (t *sqlTransaction) Begin() error {
  tx, err := t.db.Begin()
  if err != nil {
    return err
  }

  t.tx = tx
  return nil
}

func (t *sqlTransaction) Commit() error {
  return t.tx.Commit()
}

func (t *sqlTransaction) Rollback() error {
  return t.tx.Rollback()
}

func (t *sqlTransaction) Regist(c context.Context) context.Context {
  return context.WithValue(c, txKey, t.tx)
}
```

### generate proxy code

```bash
$ simple-gen-proxy -t example/transaction/service \
                   -o example/transaction/service/proxy \
                   -p proxy \
                   -n service \
                   -l github.com/ISSuh/simple-gen-proxy/example/transaction/service \
                   -x
```

### use generated proxy code

```go
type server struct {
  // single service
  foo service.Foo
  bar service.Bar

  // aggregate service
  foobar service.FooBar
}

func (s *server) init() error {
  // open db
  db, err := openDB()
  if err != nil {
    return err
  }

  // implement transaction factory function
  txFatory := func() (proxy.Transaction, error) {
    tx := repository.NewSQLTransaction(db)
    return tx, nil
  }

  // create transaction middleware
  txMiddleware := proxy.TxMiddleware(txFatory)

  // create single service
  fooRepo := repository.NewFooRepository(db)
  fooService := service.NewFooService(fooRepo)
  s.foo = proxy.NewFooProxy(fooService, txMiddleware)

  barRepo := repository.NewBarRepository(db)
  barService := service.NewBarService(barRepo)
  s.bar = proxy.NewBarProxy(barService, txMiddleware)

  // create aggregate service
  foobarService := service.NewFooBarService(s.foo, s.bar)
  s.foobar = proxy.NewFooBarProxy(foobarService, txMiddleware)
}

```

```go
  // handler example
  http.HandleFunc("POST /foo/{value}", func(w http.ResponseWriter, r *http.Request) {
    valueStr := r.PathValue("value")
    value, err := strconv.Atoi(valueStr)
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      return
    }

    d := dto.Foo{
      Value: value,
    }

    // use proxy instence
    id, err := s.foo.Create(r.Context(), d)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

    valueStr = "new foo id : " + strconv.Itoa(id)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(valueStr))
  })
```

### transaction

Save the tx instance (***sql.Tx**, ***gorm.DB**, etc...) that you use in the context through a user-defined **Transaction.Regist (ccontext.Context) error**.

If you call a service code that already has transaction middleware, determine whether there is a transaction that is already in progress through**Transaction.From(c context.Context) error**, and if the transaction already exists, participate without creating additional transactions.

```bash
# mysql log

# Transaction log on a single service call
2025-01-27T10:51:41.569361Z         3 Query     START TRANSACTION
2025-01-27T10:51:41.573225Z         3 Prepare   INSERT INTO `foo` (`value`) VALUES (?)
2025-01-27T10:51:41.573730Z         3 Execute   INSERT INTO `foo` (`value`) VALUES (123)
2025-01-27T10:51:41.574978Z         3 Close stmt
2025-01-27T10:51:41.575095Z         3 Query     COMMIT


# Transaction log on a aggreator service call
# Call two transaction proxy functions but act as one transaction
2025-01-27T10:52:09.605057Z         3 Query     START TRANSACTION
2025-01-27T10:52:09.605914Z         3 Prepare   INSERT INTO `foo` (`value`) VALUES (?)
2025-01-27T10:52:09.606331Z         3 Execute   INSERT INTO `foo` (`value`) VALUES (111)
2025-01-27T10:52:09.606894Z         3 Close stmt
2025-01-27T10:52:09.607092Z         3 Prepare   INSERT INTO `bar` (`value`) VALUES (?)
2025-01-27T10:52:09.607484Z         3 Execute   INSERT INTO `bar` (`value`) VALUES (222)
2025-01-27T10:52:09.607987Z         3 Close stmt
2025-01-27T10:52:09.608018Z         3 Query     COMMIT

# Transaction log on a aggreator service call
# Transaction rollback log due to error occurrence
2025-01-27T11:01:58.302227Z         3 Query     START TRANSACTION
2025-01-27T11:01:58.303655Z         3 Prepare   INSERT INTO `foo` (`value`) VALUES (?)
2025-01-27T11:01:58.304228Z         3 Execute   INSERT INTO `foo` (`value`) VALUES (111)
2025-01-27T11:01:58.304968Z         3 Close stmt
2025-01-27T11:01:58.305062Z         3 Query     ROLLBACK
```
