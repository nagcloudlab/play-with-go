package repository

import "transfer-service/models"

// Interface Segregation Principle - focused interface
type AccountRepository interface {
	GetAccountById(accountId string) (*models.Account, error)
	UpdateAccount(account *models.Account) error
}
