package repository

import "transfer-service/models"

type SqlAccountRepository struct {
	accounts map[string]*models.Account
}
