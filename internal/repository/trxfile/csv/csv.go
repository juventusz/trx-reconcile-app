package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"trx-reconcile-app/internal/model"
)

type Repository struct {
}

func New() *Repository {
	return &Repository{}
}

// ReadFileWithinTimeRange reads a CSV file and returns rows within the specified time range
// csv format: timestamp, id, amount, type, transactiontime
func (r *Repository) ReadFileWithinTimeRange(filename string, startTime, endTime time.Time) (map[string]model.Transaction, error) {
	// Open the CSV file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Skip the header
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Prepare a slice to store rows within the time range
	var rowsWithinRange = map[string]model.Transaction{}

	// Iterate through each row
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break // end of file reached
			}
			return nil, fmt.Errorf("failed to read row: %w", err)
		}

		// Parse the timestamp (assuming it's the first column)
		timestamp, err := time.Parse("2006-01-02 15:04:05", record[0])
		if err != nil {
			log.Printf("Skipping row with invalid timestamp: %v", err)
			continue
		}

		// Check if the timestamp is within the range
		if timestamp.Before(startTime) || timestamp.After(endTime) {
			continue
		}

		if len(record) < 4 {
			return nil, errors.New("invalid csv format")
		}

		// check trx id format
		if len(record[1]) < 36 {
			return nil, errors.New("{invalid trx id format}")
		}

		amount, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid amount , amount: %v", amount)
		}

		//construct trx object
		trx := model.Transaction{
			TransactionID:   record[1],
			Amount:          amount,
			Type:            record[3],
			TransactionTime: timestamp,
		}

		rowsWithinRange[trx.TransactionID] = trx

	}

	return rowsWithinRange, nil
}
