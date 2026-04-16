package command

// UpdateCardLimitsCmd is a command to change the spending limits of a card policy.
type UpdateCardLimitsCmd struct {
	DailyLimit  int
	WeeklyLimit int
}
