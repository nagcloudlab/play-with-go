// File: main.go
package main

import (
	"fmt"
	"transfer-service/repository"
	"transfer-service/service"
)

func main() {
	fmt.Println("Money Transfer Service v2 - SOLID Implementation")
	fmt.Println("===============================================")

	// Init / Booting phase
	//-------------------------------
	//
	// Dependency Injection - create repository once
	sqlAccountRepo := repository.GetSqlAccountRepository() // dependency

	// Inject repository into services
	upiService := service.NewUPITransferService(sqlAccountRepo)
	rtgsService := service.NewRTGSTransferService(sqlAccountRepo)

	fmt.Println("\n--- UPI Transfer Tests ---")
	demonstrateUPITransfers(upiService)

	fmt.Println("\n--- RTGS Transfer Tests ---")
	demonstrateRTGSTransfers(rtgsService)

	fmt.Println("\n--- Balance Checks ---")
	demonstrateBalanceChecks(upiService)

	fmt.Println("\n--- Error Handling Tests ---")
	demonstrateErrorHandling(upiService)
}

func demonstrateUPITransfers(service service.TransferService) {
	err := service.Transfer("1", "2", 300.00)
	if err != nil {
		fmt.Printf("UPI Transfer 1 failed: %v\n", err)
	} else {
		fmt.Println("UPI Transfer 1: Success")
	}

	err = service.Transfer("2", "3", 150.00)
	if err != nil {
		fmt.Printf("UPI Transfer 2 failed: %v\n", err)
	} else {
		fmt.Println("UPI Transfer 2: Success")
	}
}

func demonstrateRTGSTransfers(service service.TransferService) {
	// Should fail due to minimum amount requirement
	err := service.Transfer("1", "2", 100000.00)
	if err != nil {
		fmt.Printf("RTGS Transfer 1 failed: %v\n", err)
	} else {
		fmt.Println("RTGS Transfer 1: Success")
	}

	// Should succeed
	err = service.Transfer("1", "2", 250000.00)
	if err != nil {
		fmt.Printf("RTGS Transfer 2 failed: %v\n", err)
	} else {
		fmt.Println("RTGS Transfer 2: Success")
	}
}

func demonstrateBalanceChecks(service service.TransferService) {
	balance, err := service.GetAccountBalance("1")
	if err != nil {
		fmt.Printf("Balance check failed: %v\n", err)
	} else {
		fmt.Printf("Account 1 balance: %.2f\n", balance)
	}
}

func demonstrateErrorHandling(service service.TransferService) {
	// Invalid amount
	err := service.Transfer("1", "2", -100.00)
	if err != nil {
		fmt.Printf("Negative amount test: %v\n", err)
	}

	// Same account transfer
	err = service.Transfer("1", "1", 100.00)
	if err != nil {
		fmt.Printf("Same account test: %v\n", err)
	}

	// Non-existent account
	err = service.Transfer("1", "999", 100.00)
	if err != nil {
		fmt.Printf("Non-existent account test: %v\n", err)
	}

	// Insufficient balance
	err = service.Transfer("1", "2", 9999999.00)
	if err != nil {
		fmt.Printf("Insufficient balance test: %v\n", err)
	}
}
