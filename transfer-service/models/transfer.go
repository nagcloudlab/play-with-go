package models

type TransferRequest struct {
	FromAccountId string
	ToAccountId   string
	Amount        float64
	RequestId     string
}

type TransferResult struct {
	RequestId string
	Success   bool
	Error     error
}
