// File: repository/sql_account_repository.go
package repository

import (
	"context"
	"fmt"
	"sync"
	"time"
	"transfer-service/models"
)

type SqlAccountRepository struct {
	accounts map[string]*models.Account
	mutex    sync.RWMutex
}

var (
	sqlRepoInstance *SqlAccountRepository
	once            sync.Once
)

func GetSqlAccountRepository() *SqlAccountRepository {
	once.Do(func() {
		fmt.Println("[REPO] Creating SqlAccountRepository singleton with sync.Once")
		sqlRepoInstance = &SqlAccountRepository{
			accounts: initializeTestData(),
		}
		fmt.Println("[REPO] SqlAccountRepository singleton created")
	})
	return sqlRepoInstance
}

func initializeTestData() map[string]*models.Account {
	accounts := make(map[string]*models.Account)
	accounts["1"] = &models.Account{ID: "1", Name: "Alice", Balance: 1000.00}
	accounts["2"] = &models.Account{ID: "2", Name: "Bob", Balance: 500.00}
	accounts["3"] = &models.Account{ID: "3", Name: "Charlie", Balance: 750.00}
	return accounts
}

func (r *SqlAccountRepository) GetAccountById(ctx context.Context, accountId string) (*models.Account, error) {

	select {
	case <-ctx.Done():
		return nil, models.WrapContextError(ctx.Err())
	default:
	}

	time.Sleep(30 * time.Millisecond) // simulate db latency

	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if account, exists := r.accounts[accountId]; exists {
		return account, nil
	}
	return nil, models.NewAccountNotFoundError(accountId)
}

func (r *SqlAccountRepository) UpdateAccount(ctx context.Context, account *models.Account) error {

	select {
	case <-ctx.Done():
		return models.WrapContextError(ctx.Err())
	default:
	}

	time.Sleep(20 * time.Millisecond)

	r.mutex.Lock()
	defer r.mutex.Unlock()
	if _, exists := r.accounts[account.ID]; exists {
		r.accounts[account.ID] = account
		return nil
	}
	return models.NewAccountNotFoundError(account.ID)
}

func (r *SqlAccountRepository) GetMultipleAccounts(ctx context.Context, accountIds []string) ([]*models.Account, error) {

	type result struct {
		account *models.Account
		err     error
		index   int
	}

	resultChan := make(chan result, len(accountIds)) // Buffered channel to avoid goroutine leaks
	for i, id := range accountIds {
		go func(idx int, accId string) {
			acc, err := r.GetAccountById(ctx, accId)
			resultChan <- result{account: acc, err: err, index: idx}
		}(i, id)
	}

	accounts := make([]*models.Account, len(accountIds))
	var firstErr error
	for i := 0; i < len(accountIds); i++ {
		select {
		case res := <-resultChan:
			if res.err != nil && firstErr == nil {
				firstErr = res.err
			}
			accounts[res.index] = res.account
		case <-ctx.Done():
			return nil, models.WrapContextError(ctx.Err())
		}
	}
	if firstErr != nil {
		return nil, firstErr
	}
	return accounts, nil
}
