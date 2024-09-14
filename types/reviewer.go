package types

type Reviewer struct {
	Name string
	Action string // APPROVED, REQUEST_CHANGES, COMMENTED
}