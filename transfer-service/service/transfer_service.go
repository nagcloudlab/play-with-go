package service

// Interface Segregation Principle - focused contract
type TransferService interface {
	Transfer(fromAccountId, toAccountId string, amount float64) error
	GetAccountBalance(accountId string) (float64, error)
}
