// File: main.go
package main

import (
	"context"
	"fmt"
	"transfer-service/models"
	"transfer-service/repository"
	"transfer-service/service"
)

func main() {
	fmt.Println("Money Transfer Service v4 - Concurrency + Tests")
	fmt.Println("================================================")

	repo := repository.GetSqlAccountRepository()
	svc := service.NewUPITransferService(repo)

	// Single transfer demo
	err := svc.Transfer(context.Background(), "1", "2", 150)
	if err != nil {
		fmt.Printf("Transfer failed: %v\n", err)
	} else {
		fmt.Println("Transfer success")
	}

	// Bulk transfer demo
	transfers := []models.TransferRequest{
		{FromAccountId: "1", ToAccountId: "2", Amount: 10, RequestId: "REQ-1"},
		{FromAccountId: "2", ToAccountId: "3", Amount: 20, RequestId: "REQ-2"},
		{FromAccountId: "3", ToAccountId: "1", Amount: 15, RequestId: "REQ-3"},
	}
	results := svc.BulkTransfer(context.Background(), transfers)
	for _, r := range results {
		if r.Success {
			fmt.Printf("%s: SUCCESS\n", r.RequestId)
		} else {
			fmt.Printf("%s: FAILED (%v)\n", r.RequestId, r.Error)
		}
	}

	total, success := svc.GetStats()
	fmt.Printf("\nFinal Stats: total=%d, success=%d\n", total, success)
}
