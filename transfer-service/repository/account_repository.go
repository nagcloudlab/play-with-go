package repository

import (
	"context"
	"transfer-service/models"
)

type AccountRepository interface {
	GetAccountById(ctx context.Context, accountId string) (*models.Account, error)
	UpdateAccount(ctx context.Context, account *models.Account) error
	GetMultipleAccounts(ctx context.Context, accountIds []string) ([]*models.Account, error)
}
