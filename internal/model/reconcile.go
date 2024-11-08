package model

type TrxReconcile struct {
	BankName                   string                 `json:"bank"`
	TransactionProcessed       int                    `json:"trx_processed"`
	TransactionMatched         int                    `json:"trx_matched"`
	TransactionUnmatched       int                    `json:"trx_unmatched"`
	TransactionUnmatchedDetail []UnmatchedTransaction `json:"trx_unmatched_list"`
	SumAbsUnmatched            float64                `json:"sum_absolute_discrepancies"`
}

type UnmatchedTransaction struct {
	SystemTrx Transaction   `json:"system_trx"`
	BankStats BankStatement `json:"bank_stats"`
}
