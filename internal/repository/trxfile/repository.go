package trxfile

import (
	"time"
	"trx-reconcile-app/internal/model"
)

type Repository interface {
	ReadFileWithinTimeRange(filename string, startTime, endTime time.Time) (map[string]model.Transaction, error)
}
