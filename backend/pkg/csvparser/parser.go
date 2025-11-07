package csvparser

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/novanm/bank-viewer/backend/domain"
)

func Parse(fileReader io.Reader) ([]domain.Transaction, error) {
	reader := csv.NewReader(fileReader)
	reader.TrimLeadingSpace = true

	transactions := make([]domain.Transaction, 0)
	lineNumber := -1

	parseRecord := func(record []string, ln int) error {
		if len(record) != 6 {
			return fmt.Errorf("invalid format on line %d: expected 6 fields, got %d", ln, len(record))
		}

		timestamp, err := strconv.ParseInt(strings.TrimSpace(record[0]), 10, 64)
		if err != nil {
			return fmt.Errorf("invalid timestamp on line %d", ln)
		}

		name := strings.TrimSpace(record[1])
		txType := domain.TransactionType(strings.ToUpper(strings.TrimSpace(record[2])))
		if txType != domain.TypeCredit && txType != domain.TypeDebit {
			return fmt.Errorf("invalid transaction type on line %d: %s", ln, record[2])
		}

		amount, err := strconv.ParseInt(strings.TrimSpace(record[3]), 10, 64)
		if err != nil {
			return fmt.Errorf("invalid amount on line %d: %w", ln, err)
		}

		status := domain.TransactionStatus(strings.ToUpper(strings.TrimSpace(record[4])))
		if status != domain.StatusSuccess && status != domain.StatusFailed && status != domain.StatusPending {
			return fmt.Errorf("invalid status on line %d: %s", ln, record[4])
		}

		description := strings.TrimSpace(record[5])

		transactions = append(transactions, domain.Transaction{
			Timestamp:   time.Unix(timestamp, 0),
			Name:        name,
			Type:        txType,
			Amount:      amount,
			Status:      status,
			Description: description,
		})
		return nil
	}

	firstRecord, err := reader.Read()
	lineNumber++
	if err == io.EOF {
		return transactions, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read csv on line %d: %w", lineNumber, err)
	}

	isHeader := false
	if len(firstRecord) == 6 {
		expected := []string{"timestamp", "name", "type", "amount", "status", "description"}
		match := true
		for i, f := range firstRecord {
			if strings.ToLower(strings.TrimSpace(f)) != expected[i] {
				match = false
				break
			}
		}
		if match {
			isHeader = true
		}
	}

	if !isHeader {
		if err := parseRecord(firstRecord, lineNumber); err != nil {
			return nil, err
		}
	}

	for {
		lineNumber++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read csv on line %d: %w", lineNumber, err)
		}

		if err := parseRecord(record, lineNumber); err != nil {
			return nil, err
		}
	}

	return transactions, nil
}
