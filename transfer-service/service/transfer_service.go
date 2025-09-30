// File: service/transfer_service.go
package service

import (
	"context"
	"transfer-service/models"
)

type TransferService interface {
	Transfer(ctx context.Context, fromAccountId, toAccountId string, amount float64) error
	GetAccountBalance(ctx context.Context, accountId string) (float64, error)
	BulkTransfer(ctx context.Context, transfers []models.TransferRequest) []models.TransferResult
	GetStats() (int64, int64)
}
