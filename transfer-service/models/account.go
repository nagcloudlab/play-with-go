// File: models/account.go
package models

import "sync"

type Account struct {
	ID      string
	Name    string
	Balance float64
	Mutex   sync.RWMutex
}

func (a *Account) GetBalance() float64 {
	a.Mutex.RLock()
	defer a.Mutex.RUnlock()
	return a.Balance
}

func (a *Account) UpdateBalance(newBalance float64) {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()
	a.Balance = newBalance
}

func (a *Account) DebitAmount(amount float64) bool {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()
	if a.Balance >= amount {
		a.Balance -= amount
		return true
	}
	return false
}

func (a *Account) CreditAmount(amount float64) {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()
	a.Balance += amount
}
