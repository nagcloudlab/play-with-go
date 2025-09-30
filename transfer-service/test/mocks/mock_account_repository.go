// File: test/mocks/mock_account_repository.go
package mocks

import (
	"transfer-service/models"

	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) GetAccountById(accountId string) (*models.Account, error) {
	args := m.Called(accountId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) UpdateAccount(account *models.Account) error {
	args := m.Called(account)
	return args.Error(0)
}
