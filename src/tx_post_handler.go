package main

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type TransactionRequest struct {
	State         string `json:"state" validate:"required,oneof=win lost"`
	Amount        string `json:"amount" validate:"required,numeric"`
	TransactionId string `json:"transactionId" validate:"required"`
}

func TxPostHandler(ec echo.Context) error {
	c := ec.(*CustomContext)

	var transactionRequest TransactionRequest
	if err := c.Bind(&transactionRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err := c.Validator.Struct(&transactionRequest)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	sourceType := c.Request().Header.Get("Source-Type")
	if !inStringSlice(sourceType, c.Config.Sources) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Source-Type '%s' not allowed.", sourceType))
	}

	transaction := Transaction{
		State:         transactionRequest.State,
		Amount:        decimal.RequireFromString(transactionRequest.Amount),
		TransactionId: transactionRequest.TransactionId,
		Source:        sourceType,
		Canceled:      false,
	}

	if transaction.Amount.LessThanOrEqual(decimal.NewFromInt(0)) {
		return echo.NewHTTPError(http.StatusBadRequest, "Transaction Amount should by positive value.")
	}

	// Lock mutex, to secure DB from parallel access.
	c.DbMut.Lock()
	defer c.DbMut.Unlock()

	var count int
	c.Db.Model(&Transaction{}).Where("transaction_id = ?", transactionRequest.TransactionId).Count(&count)
	if count > 0 {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Tansaction with id '%s' already exists.", transactionRequest.TransactionId))
	}

	var user User
	c.Db.First(&user)

	// Double check that user exists.
	if !(user.ID > 0) {
		return echo.NewHTTPError(http.StatusInternalServerError, "User account not exist.")
	}

	if err := user.ApplyTransaction(&transaction); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Db.Create(&transaction)
	c.Db.Save(&user)

	return c.NoContent(http.StatusCreated)
}
