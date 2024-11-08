package usecase

import (
	"time"
	"trx-reconcile-app/internal/model"
)

type Usecase interface {
	ReconcileTransactionAndBanks(startTime, endTime time.Time, filenameTrxSystem string, filenameBankStatements ...string) ([]model.TrxReconcile, error)
}
