// File: repository/sql_account_repository.go
package repository

import (
	"fmt"
	"transfer-service/models"
)

// Concrete implementation of AccountRepository interface
type SqlAccountRepository struct {
	accounts map[string]*models.Account
}

// Singleton pattern for performance - shared instance
var sqlRepoInstance *SqlAccountRepository

// Factory method to get the singleton instance
func GetSqlAccountRepository() *SqlAccountRepository {
	if sqlRepoInstance == nil {
		fmt.Println("[REPO] Creating SqlAccountRepository singleton instance")
		sqlRepoInstance = &SqlAccountRepository{
			accounts: initializeTestData(),
		}
		fmt.Println("[REPO] SqlAccountRepository singleton created with test data")
	} else {
		fmt.Println("[REPO] Reusing existing SqlAccountRepository singleton")
	}
	return sqlRepoInstance
}
func initializeTestData() map[string]*models.Account {
	accounts := make(map[string]*models.Account)
	accounts["1"] = &models.Account{ID: "1", Name: "Alice", Balance: 1000.00}
	accounts["2"] = &models.Account{ID: "2", Name: "Bob", Balance: 500.00}
	accounts["3"] = &models.Account{ID: "3", Name: "Charlie", Balance: 750.00}
	return accounts
}
func (r *SqlAccountRepository) GetAccountById(accountId string) (*models.Account, error) {
	fmt.Printf("[REPO] Loading account: %s\n", accountId)

	if account, exists := r.accounts[accountId]; exists {
		fmt.Printf("[REPO] Account loaded: %s (%s), Balance: %.2f\n", account.ID, account.Name, account.Balance)
		return account, nil
	}

	fmt.Printf("[REPO] Account not found: %s\n", accountId)
	return nil, models.NewAccountNotFoundError(accountId)
}
func (r *SqlAccountRepository) UpdateAccount(account *models.Account) error {
	fmt.Printf("[REPO] Updating account: %s, New Balance: %.2f\n", account.ID, account.Balance)

	if _, exists := r.accounts[account.ID]; exists {
		r.accounts[account.ID] = account
		fmt.Println("[REPO] Account updated successfully")
		return nil
	}

	fmt.Println("[REPO] Account update failed - account not found")
	return models.NewAccountNotFoundError(account.ID)
}
