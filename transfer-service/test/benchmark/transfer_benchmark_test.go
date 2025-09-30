// File: test/benchmark/transfer_benchmark_test.go
package benchmark_test

import (
	"testing"
	"transfer-service/repository"
	"transfer-service/service"
)

func BenchmarkTransfer(b *testing.B) {
	repo := repository.GetSqlAccountRepository()
	upiService := service.NewUPITransferService(repo)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		upiService.Transfer("1", "2", 1.00)
		upiService.Transfer("2", "1", 1.00)
	}
}
