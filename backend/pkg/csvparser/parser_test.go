// bank-statement-viewer/pkg/csvparser/parser_test.go
package csvparser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse_Success(t *testing.T) {
	csvData := `1624507883, JOHN DOE, DEBIT, 250000, SUCCESS, restaurant
1624512883, COMPANY A, CREDIT, 12000000, SUCCESS, salary`
	reader := strings.NewReader(csvData)

	transactions, err := Parse(reader)

	assert.NoError(t, err)
	assert.NotNil(t, transactions)
	assert.Equal(t, 2, len(transactions))

	assert.Equal(t, int64(1624507883), transactions[0].Timestamp.Unix())
	assert.Equal(t, "JOHN DOE", transactions[0].Name)
	assert.Equal(t, int64(250000), transactions[0].Amount)
	assert.Equal(t, "restaurant", transactions[0].Description)
}

func TestParse_Error_InvalidFormat(t *testing.T) {
	csvData := `1624507883, JOHN DOE, DEBIT, 250000, SUCCESS, restaurant
1624512883, COMPANY A, CREDIT, 12000000, SUCCESS`
	reader := strings.NewReader(csvData)

	transactions, err := Parse(reader)

	assert.Error(t, err)
	assert.Nil(t, transactions)
	assert.Contains(t, err.Error(), "failed to read csv on line 1: record on line 2: wrong number of fields")
}

func TestParse_Error_InvalidAmount(t *testing.T) {
	csvData := `1624507883, JOHN DOE, DEBIT, NOT_A_NUMBER, SUCCESS, restaurant`
	reader := strings.NewReader(csvData)

	transactions, err := Parse(reader)

	assert.Error(t, err)
	assert.Nil(t, transactions)
	assert.Contains(t, err.Error(), "invalid amount on line 0: strconv.ParseInt: parsing \"NOT_A_NUMBER\": invalid syntax")
}

func generateCSVData(rows int) string {
	var sb strings.Builder
	row := "1624507883,JOHN DOE,DEBIT,250000,SUCCESS,restaurant\n"
	for i := 0; i < rows; i++ {
		sb.WriteString(row)
	}
	return sb.String()
}

func BenchmarkParse(b *testing.B) {
	csvData := generateCSVData(1000)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(csvData)

		_, err := Parse(reader)
		if err != nil {
			b.Fatal(err)
		}
	}
}
