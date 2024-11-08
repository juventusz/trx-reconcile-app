package impl

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"trx-reconcile-app/constant"
	"trx-reconcile-app/internal/model"
	"trx-reconcile-app/internal/repository/trxfile"
)

type Usecase struct {
	TrxFileRepo trxfile.Repository
}

func New(trxFileRepo trxfile.Repository) *Usecase {
	return &Usecase{
		TrxFileRepo: trxFileRepo,
	}
}

func (u *Usecase) ReconcileTransactionAndBanks(startTime, endTime time.Time, filenameTrxSystem string, filenameBankStatements ...string) ([]model.TrxReconcile, error) {
	if filenameTrxSystem == "" || len(filenameBankStatements) == 0 {
		return nil, errors.New("incomplete file paths")
	}

	trxDataRaw, err := u.TrxFileRepo.ReadFileWithinTimeRange(filenameTrxSystem, startTime, endTime)
	if err != nil {
		return nil, errors.New("error reading trx csv, err:" + err.Error())
	}

	var (
		result  []model.TrxReconcile
		trxData = trxDataRaw
	)

	for _, fileBankStatement := range filenameBankStatements {
		// Open the CSV file
		file, err := os.Open(fileBankStatement)
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

		var (
			trxReconcile                                                 model.TrxReconcile
			unmatchedTrx                                                 []model.UnmatchedTransaction
			tempUnmatchTrx                                               model.UnmatchedTransaction
			bankName                                                     = ""
			trxProcessed, trxMatched, trxUnmatched, sumAbsoluteUnmatched = 0, 0, 0, 0.0
			tempSum                                                      = 0.0
		)

		// Iterate through each row
		for {

			// part 1: Parse CSV

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

			if len(record) < 3 {
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

			trxDetail := trxData[record[1]]
			if bankName == "" {
				bankName = constant.BankCode[record[1][:4]]
			}

			// part 2: business logic
			trxProcessed++
			if trxDetail.TransactionID != "" {
				if trxDetail.Amount != amount {
					trxUnmatched++

					//sum absolute discrepancies
					tempSum = trxDetail.Amount - amount
					if tempSum < 0 {
						tempSum *= -1
					}

					tempUnmatchTrx = model.UnmatchedTransaction{
						SystemTrx: trxDetail,
						BankStats: model.BankStatement{
							UniqueID: record[1],
							Amount:   amount,
							Date:     timestamp,
						},
					}

					unmatchedTrx = append(unmatchedTrx, tempUnmatchTrx)

					sumAbsoluteUnmatched += tempSum
				} else {
					trxMatched++
				}
			}

		}

		trxReconcile = model.TrxReconcile{
			BankName:                   bankName,
			TransactionProcessed:       trxProcessed,
			TransactionMatched:         trxMatched,
			TransactionUnmatched:       trxUnmatched,
			TransactionUnmatchedDetail: unmatchedTrx,
			SumAbsUnmatched:            sumAbsoluteUnmatched,
		}

		result = append(result, trxReconcile)
	}
	return result, nil
}
