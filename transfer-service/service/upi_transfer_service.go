// File: service/upi_transfer_service.go
package service

import (
	"context"
	"fmt"
	"sync"
	"time"
	"transfer-service/models"
	"transfer-service/repository"
)

type UPITransferService struct {
	accountRepo   repository.AccountRepository
	transferCount int64
	successCount  int64
	mutex         sync.RWMutex
}

func NewUPITransferService(repo repository.AccountRepository) *UPITransferService {
	fmt.Println("[SERVICE] Creating UPITransferService with concurrency support")
	return &UPITransferService{accountRepo: repo}
}

func (s *UPITransferService) Transfer(ctx context.Context, fromId, toId string, amount float64) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	s.incrementTransferCount()
	if err := s.validateInput(fromId, toId, amount); err != nil {
		return err
	}

	accounts, err := s.accountRepo.GetMultipleAccounts(ctx, []string{fromId, toId})
	if err != nil {
		return err
	}
	from, to := accounts[0], accounts[1]
	success := s.atomicTransfer(from, to, amount)
	if !success {
		return models.NewInsufficientBalanceError(fromId, from.GetBalance(), amount)
	}

	errChan := make(chan error, 2)
	go func() { errChan <- s.accountRepo.UpdateAccount(ctx, from) }()
	go func() { errChan <- s.accountRepo.UpdateAccount(ctx, to) }()

	for i := 0; i < 2; i++ {
		select {
		case e := <-errChan:
			if e != nil {
				return e
			}
		case <-ctx.Done():
			return models.NewTimeoutError()
		}
	}

	s.incrementSuccessCount()
	return nil
}

func (s *UPITransferService) atomicTransfer(from, to *models.Account, amt float64) bool {
	var first, second *models.Account
	if from.ID < to.ID {
		first, second = from, to
	} else {
		first, second = to, from
	}
	first.Mutex.Lock()
	second.Mutex.Lock()
	defer func() {
		second.Mutex.Unlock()
		first.Mutex.Unlock()
	}()
	if from.Balance >= amt {
		from.Balance -= amt
		to.Balance += amt
		return true
	}
	return false
}

func (s *UPITransferService) validateInput(from, to string, amt float64) error {
	if amt <= 0 {
		return models.NewInvalidAmountError(amt)
	}
	if from == to {
		return &models.TransferError{Code: "SAME_ACCOUNT_TRANSFER", Message: "Cannot transfer to same account"}
	}
	return nil
}

func (s *UPITransferService) GetAccountBalance(ctx context.Context, accountId string) (float64, error) {
	acc, err := s.accountRepo.GetAccountById(ctx, accountId)
	if err != nil {
		return 0, err
	}
	return acc.GetBalance(), nil
}

func (s *UPITransferService) BulkTransfer(ctx context.Context, transfers []models.TransferRequest) []models.TransferResult {
	const workers = 3
	jobChan := make(chan models.TransferRequest, len(transfers))
	resChan := make(chan models.TransferResult, len(transfers))

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for tr := range jobChan {
				err := s.Transfer(ctx, tr.FromAccountId, tr.ToAccountId, tr.Amount)
				resChan <- models.TransferResult{RequestId: tr.RequestId, Success: err == nil, Error: err}
			}
		}(i)
	}

	for _, t := range transfers {
		jobChan <- t
	}
	close(jobChan)

	go func() {
		wg.Wait()
		close(resChan)
	}()

	var results []models.TransferResult
	for r := range resChan {
		results = append(results, r)
	}
	return results
}
func (s *UPITransferService) incrementTransferCount() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.transferCount++
}
func (s *UPITransferService) incrementSuccessCount() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.successCount++
}
func (s *UPITransferService) GetStats() (int64, int64) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.transferCount, s.successCount
}
