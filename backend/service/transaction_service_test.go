// bank-statement-viewer/service/transaction_service_test.go
package service

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/novanm/bank-viewer/backend/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Store(ctx context.Context, txs []domain.Transaction) error {
	args := m.Called(ctx, txs)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetAll(ctx context.Context) ([]domain.Transaction, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Transaction), args.Error(1)
}

var t1 = time.Now()
var t2 = t1.Add(1 * time.Hour)
var t3 = t1.Add(2 * time.Hour)

var mockData = []domain.Transaction{

	{Timestamp: t1, Name: "COMPANY A", Type: domain.TypeCredit, Amount: 1000, Status: domain.StatusSuccess},
	{Timestamp: t2, Name: "RESTAURANT", Type: domain.TypeDebit, Amount: 100, Status: domain.StatusSuccess},

	{Timestamp: t2, Name: "E-COMMERCE", Type: domain.TypeDebit, Amount: 50, Status: domain.StatusFailed},

	{Timestamp: t3, Name: "TRANSFER", Type: domain.TypeCredit, Amount: 200, Status: domain.StatusPending},
}

func TestGetBalance_Success(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	mockRepo.On("GetAll", mock.Anything).Return(mockData, nil)

	s := NewTransactionService(mockRepo)

	balance, err := s.GetBalance(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, balance)

	assert.Equal(t, int64(900), balance.TotalBalance)
	mockRepo.AssertExpectations(t)
}

func TestGetIssues_PaginationAndSorting(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	mockRepo.On("GetAll", mock.Anything).Return(mockData, nil)
	s := NewTransactionService(mockRepo)

	params := domain.PaginationParams{
		Page:    1,
		Limit:   10,
		SortBy:  "amount",
		SortDir: "asc",
	}

	issues, err := s.GetIssues(context.Background(), params)

	assert.NoError(t, err)
	assert.NotNil(t, issues)

	assert.Equal(t, 2, len(issues.Transactions))
	assert.Equal(t, 2, issues.Metadata.TotalItems)

	assert.Equal(t, int64(50), issues.Transactions[0].Amount)
	assert.Equal(t, int64(200), issues.Transactions[1].Amount)

	params2 := domain.PaginationParams{
		Page:    2,
		Limit:   1,
		SortBy:  "amount",
		SortDir: "asc",
	}

	issues2, err2 := s.GetIssues(context.Background(), params2)

	assert.NoError(t, err2)
	assert.NotNil(t, issues2)
	assert.Equal(t, 2, issues2.Metadata.TotalItems)             // Total tetap 2
	assert.Equal(t, 2, issues2.Metadata.TotalPages)             // Sekarang ada 2 halaman
	assert.Equal(t, 2, issues2.Metadata.CurrentPage)            // Kita di halaman 2
	assert.Equal(t, 1, len(issues2.Transactions))               // Hanya 1 item
	assert.Equal(t, int64(200), issues2.Transactions[0].Amount) // Item kedua

	mockRepo.AssertExpectations(t)
}

func TestProcessUpload_Success(t *testing.T) {

	csvData := `1624507883, JOHN DOE, DEBIT, 25000, SUCCESS, restaurant`
	reader := strings.NewReader(csvData)

	mockRepo := new(MockTransactionRepository)
	mockRepo.On("Store", mock.Anything, mock.AnythingOfType("[]domain.Transaction")).Return(nil)

	s := NewTransactionService(mockRepo)

	err := s.ProcessUpload(context.Background(), reader)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Store", mock.Anything, mock.AnythingOfType("[]domain.Transaction"))
}

func TestProcessUpload_ParseError(t *testing.T) {
	csvData := `1624507883, JOHN DOE, DEBIT`
	reader := strings.NewReader(csvData)

	mockRepo := new(MockTransactionRepository)

	s := NewTransactionService(mockRepo)

	err := s.ProcessUpload(context.Background(), reader)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid format")

	mockRepo.AssertNotCalled(t, "Store", mock.Anything, mock.Anything)
}

func generateMockData(rows int) []domain.Transaction {
	data := make([]domain.Transaction, 0, rows)
	for i := 0; i < rows; i++ {
		status := domain.StatusSuccess
		if i%3 == 0 {
			status = domain.StatusFailed
		} else if i%5 == 0 {
			status = domain.StatusPending
		}

		data = append(data, domain.Transaction{
			Timestamp: time.Now().Add(time.Duration(i) * time.Second),
			Name:      "E-COMMERCE " + strconv.Itoa(i),
			Type:      domain.TypeDebit,
			Amount:    int64(i * 100), // Buat amount berbeda agar sorting bekerja
			Status:    status,
		})
	}
	return data
}

func BenchmarkGetIssues(b *testing.B) {
	largeMockData := generateMockData(10000)

	mockRepo := new(MockTransactionRepository)
	mockRepo.On("GetAll", mock.Anything).Return(largeMockData, nil)

	s := NewTransactionService(mockRepo)

	params := domain.PaginationParams{
		Page:    1,
		Limit:   100,
		SortBy:  "amount",
		SortDir: "desc",
	}
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := s.GetIssues(ctx, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}
