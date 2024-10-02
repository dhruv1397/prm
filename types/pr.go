package types

import "strings"

type PullRequest struct {
	Number           int      `json:"number" yaml:"number"`
	Title            string   `json:"title" yaml:"title"`
	SCMProviderType  string   `json:"scm_provider_type" yaml:"scm_provider_type"`
	SCMProviderName  string   `json:"scm_provider_name" yaml:"scm_provider_name"`
	URL              string   `json:"url" yaml:"url"`
	State            string   `json:"state" yaml:"state"`
	Approved         []string `json:"approved" yaml:"approved"`
	Commented        []string `json:"commented" yaml:"commented"`
	RequestedChanges []string `json:"requested_changes" yaml:"requested_changes"`
	Mergeable        string   `json:"mergeable" yaml:"mergeable"`
}

type PrintablePullRequest struct {
	NumberRaw          int
	SCMProviderTypeRaw string
	SCMProviderNameRaw string
	Number             []string
	Title              []string
	SCMProviderName    []string
	URL                []string
	State              []string
	Approved           []string
	Commented          []string
	RequestedChanges   []string
	Mergeable          []string
	MaxRows            int
}

type PullRequestResponse struct {
	PR          *PullRequest
	PrintablePR *PrintablePullRequest
}

func ComparePullRequest(a, b *PullRequest) int {
	if a.SCMProviderType != b.SCMProviderType {
		return strings.Compare(a.SCMProviderType, b.SCMProviderType)
	}
	if a.SCMProviderName != b.SCMProviderName {
		return strings.Compare(a.SCMProviderName, b.SCMProviderName)
	}
	if a.Number > b.Number {
		return -1
	}
	if a.Number < b.Number {
		return 1
	}
	return 0
}

func ComparePrintablePullRequest(a, b *PrintablePullRequest) int {
	if a.SCMProviderTypeRaw != b.SCMProviderTypeRaw {
		return strings.Compare(a.SCMProviderTypeRaw, b.SCMProviderTypeRaw)
	}
	if a.SCMProviderNameRaw != b.SCMProviderNameRaw {
		return strings.Compare(a.SCMProviderNameRaw, b.SCMProviderNameRaw)
	}
	if a.NumberRaw > b.NumberRaw {
		return -1
	}
	if a.NumberRaw < b.NumberRaw {
		return 1
	}
	return 0
}
