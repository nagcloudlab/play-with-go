// File: service/upi_transfer_service.go
package service

import (
	"fmt"
	"transfer-service/models"
	"transfer-service/repository"
)

// UPI implementation of TransferService interface
type UPITransferService struct {
	accountRepo repository.AccountRepository // Dependency Inversion - depends on interface
}

// Constructor with Dependency Injection
func NewUPITransferService(accountRepo repository.AccountRepository) *UPITransferService {
	fmt.Println("[SERVICE] Creating UPITransferService with injected dependencies")
	return &UPITransferService{
		accountRepo: accountRepo,
	}
}

// Single Responsibility - only handles transfer logic
func (s *UPITransferService) Transfer(fromAccountId, toAccountId string, amount float64) error {
	fmt.Printf("[SERVICE] Starting UPI transfer: %.2f from %s to %s\n", amount, fromAccountId, toAccountId)

	// Input validation
	if err := s.validateTransferInput(fromAccountId, toAccountId, amount); err != nil {
		return err
	}

	// Load accounts
	fromAccount, err := s.accountRepo.GetAccountById(fromAccountId)
	if err != nil {
		return err
	}

	toAccount, err := s.accountRepo.GetAccountById(toAccountId)
	if err != nil {
		return err
	}

	// Business validation
	if fromAccount.Balance < amount {
		return models.NewInsufficientBalanceError(fromAccountId, fromAccount.Balance, amount)
	}

	// Perform transfer
	fmt.Printf("[SERVICE] Debiting %.2f from account %s via UPI\n", amount, fromAccountId)
	fromAccount.Balance = fromAccount.Balance - amount

	fmt.Printf("[SERVICE] Crediting %.2f to account %s via UPI\n", amount, toAccountId)
	toAccount.Balance = toAccount.Balance + amount

	// Update accounts
	if err := s.accountRepo.UpdateAccount(fromAccount); err != nil {
		return err
	}

	if err := s.accountRepo.UpdateAccount(toAccount); err != nil {
		return err
	}

	fmt.Println("[SERVICE] UPI Transfer completed successfully")
	return nil
}

// Single Responsibility - focused validation
func (s *UPITransferService) validateTransferInput(fromAccountId, toAccountId string, amount float64) error {
	if amount <= 0 {
		return models.NewInvalidAmountError(amount)
	}

	if fromAccountId == toAccountId {
		return &models.TransferError{
			Code:    "SAME_ACCOUNT_TRANSFER",
			Message: "Cannot transfer to the same account",
			Details: map[string]interface{}{"accountId": fromAccountId},
		}
	}

	if fromAccountId == "" || toAccountId == "" {
		return &models.TransferError{
			Code:    "EMPTY_ACCOUNT_ID",
			Message: "Account IDs cannot be empty",
			Details: map[string]interface{}{
				"fromAccountId": fromAccountId,
				"toAccountId":   toAccountId,
			},
		}
	}

	return nil
}

// Single Responsibility - only handles balance inquiry
func (s *UPITransferService) GetAccountBalance(accountId string) (float64, error) {
	fmt.Printf("[SERVICE] Getting balance for account: %s\n", accountId)

	account, err := s.accountRepo.GetAccountById(accountId)
	if err != nil {
		return 0, err
	}

	return account.Balance, nil
}
