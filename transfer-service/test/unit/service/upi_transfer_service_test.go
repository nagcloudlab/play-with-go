// File: test/unit/service/upi_transfer_service_test.go
package service_test

import (
	"testing"
	"transfer-service/models"
	"transfer-service/service"
	"transfer-service/test/helpers"
	"transfer-service/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUPITransferService_Transfer_Success(t *testing.T) {

	mockRepo := new(mocks.MockAccountRepository)
	upiService := service.NewUPITransferService(mockRepo)

	fromAccount := helpers.CreateTestAccount("1", "Alice", 1000.00)
	toAccount := helpers.CreateTestAccount("2", "Bob", 500.00)

	mockRepo.On("GetAccountById", "1").Return(fromAccount, nil)
	mockRepo.On("GetAccountById", "2").Return(toAccount, nil)
	mockRepo.On("UpdateAccount", mock.AnythingOfType("*models.Account")).Return(nil).Twice()

	err := upiService.Transfer("1", "2", 300.00)

	assert.NoError(t, err)
	assert.Equal(t, 700.00, fromAccount.Balance)
	assert.Equal(t, 800.00, toAccount.Balance)
	mockRepo.AssertExpectations(t)
}

func TestUPITransferService_Transfer_InsufficientBalance(t *testing.T) {
	mockRepo := new(mocks.MockAccountRepository)
	upiService := service.NewUPITransferService(mockRepo)

	fromAccount := helpers.CreateTestAccount("1", "Alice", 100.00)
	toAccount := helpers.CreateTestAccount("2", "Bob", 500.00)

	mockRepo.On("GetAccountById", "1").Return(fromAccount, nil)
	mockRepo.On("GetAccountById", "2").Return(toAccount, nil)

	err := upiService.Transfer("1", "2", 300.00)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "INSUFFICIENT_BALANCE")
	mockRepo.AssertExpectations(t)
}

func TestUPITransferService_Transfer_AccountNotFound(t *testing.T) {
	mockRepo := new(mocks.MockAccountRepository)
	upiService := service.NewUPITransferService(mockRepo)

	mockRepo.On("GetAccountById", "999").Return(nil, models.NewAccountNotFoundError("999"))

	err := upiService.Transfer("999", "2", 300.00)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ACCOUNT_NOT_FOUND")
	mockRepo.AssertExpectations(t)
}

func TestUPITransferService_Transfer_InvalidInputs(t *testing.T) {

	testCases := []struct {
		name          string
		fromAccountId string
		toAccountId   string
		amount        float64
		expectedError string
	}{
		{"Negative amount", "1", "2", -100.00, "INVALID_AMOUNT"},
		{"Zero amount", "1", "2", 0.00, "INVALID_AMOUNT"},
		{"Same account", "1", "1", 100.00, "SAME_ACCOUNT_TRANSFER"},
		{"Empty from ID", "", "2", 100.00, "EMPTY_ACCOUNT_ID"},
		{"Empty to ID", "1", "", 100.00, "EMPTY_ACCOUNT_ID"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAccountRepository)
			upiService := service.NewUPITransferService(mockRepo)
			err := upiService.Transfer(tc.fromAccountId, tc.toAccountId, tc.amount)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedError)
		})
	}
}

func TestUPITransferService_GetAccountBalance_Success(t *testing.T) {
	mockRepo := new(mocks.MockAccountRepository)
	upiService := service.NewUPITransferService(mockRepo)

	account := helpers.CreateTestAccount("1", "Alice", 1000.00)
	mockRepo.On("GetAccountById", "1").Return(account, nil)

	balance, err := upiService.GetAccountBalance("1")

	assert.NoError(t, err)
	assert.Equal(t, 1000.00, balance)
	mockRepo.AssertExpectations(t)
}
