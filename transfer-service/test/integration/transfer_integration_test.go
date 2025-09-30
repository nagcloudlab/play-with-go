// File: test/integration/transfer_integration_test.go
package integration_test

import (
	"testing"
	"transfer-service/repository"
	"transfer-service/service"

	"github.com/stretchr/testify/assert"
)

func TestTransferIntegration_SuccessfulTransfer(t *testing.T) {
	repo := repository.GetSqlAccountRepository()
	upiService := service.NewUPITransferService(repo)

	initialBalance1, _ := upiService.GetAccountBalance("1")
	initialBalance2, _ := upiService.GetAccountBalance("2")

	transferAmount := 200.00
	err := upiService.Transfer("1", "2", transferAmount)

	assert.NoError(t, err)

	finalBalance1, _ := upiService.GetAccountBalance("1")
	finalBalance2, _ := upiService.GetAccountBalance("2")

	assert.Equal(t, initialBalance1-transferAmount, finalBalance1)
	assert.Equal(t, initialBalance2+transferAmount, finalBalance2)
}
