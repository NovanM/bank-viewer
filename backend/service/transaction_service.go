package service

import (
	"context"
	"io"
	"math"
	"sort"

	"github.com/novanm/bank-viewer/backend/domain"
	"github.com/novanm/bank-viewer/backend/pkg/csvparser"
)

type TransactionService struct {
	repo domain.TransactionRepository
}

func NewTransactionService(repo domain.TransactionRepository) *TransactionService {
	return &TransactionService{
		repo: repo,
	}
}

func (s *TransactionService) ProcessUpload(ctx context.Context, fileReader io.Reader) error {
	transactions, err := csvparser.Parse(fileReader)
	if err != nil {
		return err
	}

	err = s.repo.Store(ctx, transactions)
	if err != nil {
		return err
	}

	return nil
}

func (s *TransactionService) GetBalance(ctx context.Context) (*domain.BalanceResponse, error) {
	transactions, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var totalBalance int64
	for _, tx := range transactions {
		if tx.Status == domain.StatusSuccess {
			switch tx.Type {
			case domain.TypeCredit:
				totalBalance += tx.Amount
			case domain.TypeDebit:
				totalBalance -= tx.Amount
			}
		}
	}

	return &domain.BalanceResponse{
		TotalBalance: totalBalance,
	}, nil
}

func (s *TransactionService) GetIssues(ctx context.Context, params domain.PaginationParams) (*domain.IssuesResponse, error) {

	transactions, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	issues := make([]domain.Transaction, 0)
	for _, tx := range transactions {
		if tx.Status == domain.StatusFailed || tx.Status == domain.StatusPending {
			issues = append(issues, tx)
		}
	}

	sort.Slice(issues, func(i, j int) bool {
		switch params.SortBy {
		case "amount":
			if params.SortDir == "asc" {
				return issues[i].Amount < issues[j].Amount
			}
			return issues[i].Amount > issues[j].Amount
		case "name":
			if params.SortDir == "asc" {
				return issues[i].Name < issues[j].Name
			}
			return issues[i].Name > issues[j].Name
		default:
			if params.SortDir == "asc" {
				return issues[i].Timestamp.Before(issues[j].Timestamp)
			}
			return issues[i].Timestamp.After(issues[j].Timestamp)
		}
	})

	totalItems := len(issues)
	totalPages := int(math.Ceil(float64(totalItems) / float64(params.Limit)))

	startIndex := (params.Page - 1) * params.Limit
	endIndex := startIndex + params.Limit

	if startIndex > totalItems {
		startIndex = totalItems
	}
	if endIndex > totalItems {
		endIndex = totalItems
	}

	pageData := issues[startIndex:endIndex]

	metadata := domain.PaginationMetadata{
		CurrentPage: params.Page,
		PageSize:    params.Limit,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
	}

	response := &domain.IssuesResponse{
		Transactions: pageData,
		Metadata:     metadata,
	}

	return response, nil
}
