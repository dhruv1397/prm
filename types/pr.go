package types

type PullRequest struct {
	SCMProviderType string
	URL string
	State string
	Reviewers []*Reviewer
}