// File: test/helpers/test_helpers.go
package helpers

import "transfer-service/models"

func CreateTestAccount(id, name string, balance float64) *models.Account {
	return &models.Account{
		ID:      id,
		Name:    name,
		Balance: balance,
	}
}

func CreateTestAccounts() []*models.Account {
	return []*models.Account{
		CreateTestAccount("1", "Alice", 1000.00),
		CreateTestAccount("2", "Bob", 500.00),
		CreateTestAccount("3", "Charlie", 750.00),
	}
}
