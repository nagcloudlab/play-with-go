package models

import "fmt"

// Custom error types for better error handling
type TransferError struct {
	Code    string
	Message string
	Details map[string]interface{}
}

func (e *TransferError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Predefined error types
func NewAccountNotFoundError(accountId string) *TransferError {
	return &TransferError{
		Code:    "ACCOUNT_NOT_FOUND",
		Message: fmt.Sprintf("Account %s not found", accountId),
		Details: map[string]interface{}{"accountId": accountId},
	}
}

func NewInsufficientBalanceError(accountId string, balance, amount float64) *TransferError {
	return &TransferError{
		Code:    "INSUFFICIENT_BALANCE",
		Message: fmt.Sprintf("Account %s has insufficient balance", accountId),
		Details: map[string]interface{}{
			"accountId": accountId,
			"balance":   balance,
			"required":  amount,
		},
	}
}

func NewInvalidAmountError(amount float64) *TransferError {
	return &TransferError{
		Code:    "INVALID_AMOUNT",
		Message: fmt.Sprintf("Invalid transfer amount: %.2f", amount),
		Details: map[string]interface{}{"amount": amount},
	}
}
