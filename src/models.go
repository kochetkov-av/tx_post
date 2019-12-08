package main

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	gorm.Model
	State         string
	Amount        decimal.Decimal `gorm:"type:decimal"`
	TransactionId string          `gorm:"unique_index"`
	Source        string
	Canceled      bool `gorm:"NOT NULL"`
}

type User struct {
	gorm.Model
	Balance decimal.Decimal `gorm:"type:decimal"`
}

func (user *User) ApplyTransaction(transaction *Transaction) error {
	var newUserBalance decimal.Decimal
	if transaction.State == "win" {
		newUserBalance = user.Balance.Add(transaction.Amount)
	} else if transaction.State == "lost" {
		newUserBalance = user.Balance.Sub(transaction.Amount)
	} else {
		return errors.New("not supported transaction state")
	}

	if newUserBalance.LessThan(decimal.NewFromInt(0)) {
		return errors.New("failed to process transaction, negative balance")
	}

	user.Balance = newUserBalance
	return nil
}

func (user *User) CancelTransaction(transaction *Transaction) error {
	transaction.Canceled = true

	revertTransaction := *transaction

	revertTransaction.Amount = revertTransaction.Amount.Mul(decimal.NewFromInt(-1))

	return user.ApplyTransaction(&revertTransaction)
}
