// bank-statement-viewer/repository/memory/transaction_repo_test.go
package memory

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/novanm/bank-viewer/backend/domain"
	"github.com/stretchr/testify/assert"
)

func TestStoreAndGetAll_Success(t *testing.T) {

	repo := NewMemoryRepository()
	testData := []domain.Transaction{
		{Name: "Test 1", Amount: 100},
		{Name: "Test 2", Amount: 200},
	}
	ctx := context.Background()

	err := repo.Store(ctx, testData)
	assert.NoError(t, err)

	data, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(data))
	assert.Equal(t, "Test 1", data[0].Name)

	testData2 := []domain.Transaction{
		{Name: "Test 3", Amount: 300},
	}
	err = repo.Store(ctx, testData2)
	assert.NoError(t, err)

	data2, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(data2))
	assert.Equal(t, "Test 3", data2[0].Name)
}

func TestRepository_Concurrency(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	initialData := []domain.Transaction{{Name: "Initial", Amount: 1}}
	repo.Store(ctx, initialData)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		time.Sleep(50 * time.Millisecond) // Biarkan pembaca jalan dulu
		newData := []domain.Transaction{{Name: "New Data", Amount: 999}}
		repo.Store(ctx, newData)
	}()

	go func() {
		defer wg.Done()
		data, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, data)
	}()

	go func() {
		defer wg.Done()
		data, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, data)
	}()

	wg.Wait()

	finalData, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(finalData))
	assert.Equal(t, "New Data", finalData[0].Name)

}
