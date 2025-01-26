package proxy

import (
	"context"
	"errors"
)

var (
	errNilTransactionFactory = errors.New("transaction factory is nil")
	errRollbackTransaction   = errors.New("rollback transaction")
)

// Transaction is an interface that defines the methods to manage transactions.
type Transaction interface {
	// Begin begins the transaction.
	Begin() error

	// Commit commits the transaction.
	Commit() error

	// Rollback rolls back the transaction.
	Rollback() error

	// Regist and returns a context with the transaction.
	Regist(c context.Context) context.Context

	// From sets the transaction from the context.
	From(c context.Context) error
}

// NewTransaction is a function that creates a new transaction.
type TransactionFactory func() (Transaction, error)

// TxMiddleware is a function that returns a middleware that manages transactions.
func TxMiddleware(creator TransactionFactory) func(func(c context.Context) error) func(context.Context) error {
	return func(next func(c context.Context) error) func(context.Context) error {
		return func(c context.Context) error {
			if creator == nil {
				return errNilTransactionFactory
			}

			tx, err := creator()
			if err != nil {
				return err
			}

			err = tx.From(c)
			if err == nil {
				return subTransaction(c, next, tx)
			}

			newTransaction(c, next, tx)
			return nil
		}
	}
}

// newTransaction creates a new transaction.
func newTransaction(
	c context.Context, next func(c context.Context) error, tx Transaction,
) error {
	err := tx.Begin()
	if err != nil {
		return err
	}

	c = tx.Regist(c)
	if err := next(c); err != nil {
		tx.Rollback()
		return errors.Join(errRollbackTransaction, err)
	}

	tx.Commit()
	return nil
}

// subTransaction manages the transaction.
func subTransaction(
	c context.Context, next func(c context.Context) error, tx Transaction,
) error {
	if err := next(c); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return errors.Join(errRollbackTransaction, err)
	}
	return nil
}
