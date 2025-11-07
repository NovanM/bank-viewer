package domain

import (
	"context"
	"io"
)

type TransactionService interface {
	ProcessUpload(ctx context.Context, fileReader io.Reader) error
	GetBalance(ctx context.Context) (*BalanceResponse, error)
	GetIssues(ctx context.Context, params PaginationParams) (*IssuesResponse, error)
}

type TransactionRepository interface {
	Store(ctx context.Context, transactions []Transaction) error
	GetAll(ctx context.Context) ([]Transaction, error)
}
