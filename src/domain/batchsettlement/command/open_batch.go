package command

// OpenBatchCmd is the command to initiate a new daily batch settlement cycle.
type OpenBatchCmd struct {
	SettlementDate   string
	OperationalRegion string
}
