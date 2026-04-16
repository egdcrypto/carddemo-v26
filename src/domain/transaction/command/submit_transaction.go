package command

// SubmitTransactionCmd represents a command to submit a new financial transaction.
type SubmitTransactionCmd struct {
	TransactionID string
	AccountID     string
	CardID        string
	Amount        float64
	Type          string
	AccountStatus string
}
