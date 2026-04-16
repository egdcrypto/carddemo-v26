package dto

// QueryParams represents supported query parameters for listing transactions.
type QueryParams struct {
	StartDate string
	EndDate   string
	Status    string
	Page      int
	Limit     int
}
