package memory

import (
	"context"
	"sync"

	"github.com/novanm/bank-viewer/backend/domain"
)

type memoryRepository struct {
	transactions []domain.Transaction
	mu           sync.RWMutex
}

func NewMemoryRepository() domain.TransactionRepository {
	return &memoryRepository{
		transactions: make([]domain.Transaction, 0),
	}
}

func (m *memoryRepository) GetAll(ctx context.Context) ([]domain.Transaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	transactionsCopy := make([]domain.Transaction, len(m.transactions))
	copy(transactionsCopy, m.transactions)
	return transactionsCopy, nil
}

func (m *memoryRepository) Store(ctx context.Context, transactions []domain.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.transactions = transactions
	return nil
}
