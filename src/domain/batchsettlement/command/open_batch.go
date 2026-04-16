package command

// OpenBatchCmd is the command to open a new settlement batch.
type OpenBatchCmd struct {
	BatchID    string
	MerchantID string
	Currency   string
}
