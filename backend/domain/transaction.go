package domain

import "time"

type TransactionStatus string

const (
	StatusPending TransactionStatus = "PENDING"
	StatusSuccess TransactionStatus = "SUCCESS"
	StatusFailed  TransactionStatus = "FAILED"
)

type TransactionType string

const (
	TypeCredit  TransactionType = "CREDIT"
	TypeDebit   TransactionType = "DEBIT"
	TypeUnknown TransactionType = ""
)

type Transaction struct {
	Timestamp   time.Time         `json:"timestamp"`
	Name        string            `json:"name"`
	Type        TransactionType   `json:"type"`
	Amount      int64             `json:"amount"`
	Status      TransactionStatus `json:"status"`
	Description string            `json:"description"`
}

type BalanceResponse struct {
	TotalBalance int64 `json:"total_balance"`
}

type PaginationParams struct {
	Page    int
	Limit   int
	SortBy  string
	SortDir string
}

type PaginationMetadata struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
}

type IssuesResponse struct {
	Transactions []Transaction      `json:"transactions"`
	Metadata     PaginationMetadata `json:"metadata"`
}
