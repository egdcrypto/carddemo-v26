package command

// OpenAccountCmd is a command to open a new account.
type OpenAccountCmd struct {
	UserProfileID  string
	InitialStatus  string
	AccountType    string
	CurrentBalance float64 // Used to test invariants about closure/zero balance
	IsClosed       bool    // Used to test invariants about irreversibility
}
