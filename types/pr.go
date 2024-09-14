package types

import "strings"

type PullRequest struct {
	Number           int
	Title            string
	SCMProviderType  string
	URL              string
	State            string
	Approved         []string
	Commented        []string
	RequestedChanges []string
	Mergeable        string
}

type PrintablePullRequest struct {
	NumberRaw          int
	SCMProviderTypeRaw string
	Number             []string
	Title              []string
	SCMProviderType    []string
	URL                []string
	State              []string
	Approved           []string
	Commented          []string
	RequestedChanges   []string
	Mergeable          []string
	MaxRows            int
}

func ComparePrintablePullRequest(a, b *PrintablePullRequest) int {
	if a.SCMProviderTypeRaw == b.SCMProviderTypeRaw {
		if a.NumberRaw > b.NumberRaw {
			return 1
		} else if a.NumberRaw == b.NumberRaw {
			return 0
		}
		return -1
	}

	return strings.Compare(a.SCMProviderTypeRaw, b.SCMProviderTypeRaw)
}
