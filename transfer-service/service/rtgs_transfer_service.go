// File: service/rtgs_transfer_service.go
package service

import (
	"fmt"
	"transfer-service/models"
	"transfer-service/repository"
)

// Open/Closed Principle - new implementation without modifying existing code
type RTGSTransferService struct {
	accountRepo repository.AccountRepository
}

func NewRTGSTransferService(accountRepo repository.AccountRepository) *RTGSTransferService {
	fmt.Println("[SERVICE] Creating RTGSTransferService with injected dependencies")
	return &RTGSTransferService{
		accountRepo: accountRepo,
	}
}

func (s *RTGSTransferService) Transfer(fromAccountId, toAccountId string, amount float64) error {
	fmt.Printf("[SERVICE] Starting RTGS transfer: %.2f from %s to %s\n", amount, fromAccountId, toAccountId)

	// RTGS specific validation (minimum amount)
	if amount < 200000 {
		return &models.TransferError{
			Code:    "RTGS_MINIMUM_AMOUNT",
			Message: "RTGS transfers require minimum amount of 200,000",
			Details: map[string]interface{}{"amount": amount, "minimum": 200000},
		}
	}

	// Reuse common validation
	if err := s.validateCommonInput(fromAccountId, toAccountId, amount); err != nil {
		return err
	}

	// Core logic
	fromAccount, err := s.accountRepo.GetAccountById(fromAccountId)
	if err != nil {
		return err
	}

	toAccount, err := s.accountRepo.GetAccountById(toAccountId)
	if err != nil {
		return err
	}

	if fromAccount.Balance < amount {
		return models.NewInsufficientBalanceError(fromAccountId, fromAccount.Balance, amount)
	}

	fmt.Printf("[SERVICE] Debiting %.2f from account %s via RTGS\n", amount, fromAccountId)
	fromAccount.Balance = fromAccount.Balance - amount

	fmt.Printf("[SERVICE] Crediting %.2f to account %s via RTGS\n", amount, toAccountId)
	toAccount.Balance = toAccount.Balance + amount

	if err := s.accountRepo.UpdateAccount(fromAccount); err != nil {
		return err
	}

	if err := s.accountRepo.UpdateAccount(toAccount); err != nil {
		return err
	}

	fmt.Println("[SERVICE] RTGS Transfer completed successfully")
	return nil
}

func (s *RTGSTransferService) validateCommonInput(fromAccountId, toAccountId string, amount float64) error {
	if amount <= 0 {
		return models.NewInvalidAmountError(amount)
	}
	if fromAccountId == toAccountId {
		return &models.TransferError{
			Code:    "SAME_ACCOUNT_TRANSFER",
			Message: "Cannot transfer to the same account",
		}
	}
	return nil
}

func (s *RTGSTransferService) GetAccountBalance(accountId string) (float64, error) {
	account, err := s.accountRepo.GetAccountById(accountId)
	if err != nil {
		return 0, err
	}
	return account.Balance, nil
}
