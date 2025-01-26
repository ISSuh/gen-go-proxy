package proxy

// import (
// 	"context"
// 	"errors"
// )

// // Transaction is an interface that defines the methods to manage transactions.
// type Transaction interface {
// 	// Begin begins the transaction.
// 	Begin() error

// 	// Commit commits the transaction.
// 	Commit() error

// 	// Rollback rolls back the transaction.
// 	Rollback() error

// 	// Regist and returns a context with the transaction.
// 	Regist(c context.Context) context.Context

// 	// From sets the transaction from the context.
// 	From(c context.Context) error
// }

// // NewTransaction is a function that creates a new transaction.
// type TransactionFactory func() (Transaction, error)

// // TxHepler is a helper struct that manages transactions.
// type TxHepler struct {
// 	creator TransactionFactory
// }

// func NewTxHepler(creator TransactionFactory) TxHepler {
// 	return TxHepler{
// 		creator: creator,
// 	}
// }

// // Transaction is a function that returns a user transaction function.
// func (t TxHepler) Transaction() func(c context.Context, f func(c context.Context) error) error {
// 	return func(c context.Context, f func(c context.Context) error) error {
// 		if t.creator == nil {
// 			return errors.New("transaction  found")
// 		}

// 		tx, err := t.creator()
// 		if err != nil {
// 			return err
// 		}

// 		err = tx.From(c)
// 		if err == nil {
// 			return t.subTransaction(c, f, tx)
// 		}

// 		t.newTransaction(c, f, tx)
// 		return nil
// 	}
// }

// func (t TxHepler) newTransaction(
// 	c context.Context, f func(c context.Context) error, tx Transaction,
// ) error {
// 	err := tx.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	c = tx.Regist(c)
// 	if err := f(c); err != nil {
// 		tx.Rollback()
// 		return errors.Join(errors.New("rollback transaction"), err)
// 	}

// 	tx.Commit()
// 	return nil
// }

// func (t TxHepler) subTransaction(
// 	c context.Context, f func(c context.Context) error, tx Transaction,
// ) error {
// 	err := f(c)
// 	if err != nil {
// 		if err := tx.Rollback(); err != nil {
// 			return err
// 		}

// 		return errors.Join(errors.New("rollback transaction"), err)
// 	}
// 	return nil
// }
