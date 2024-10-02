package types

type EmailData struct {
	Email string `json:"email"`
}

type EmailResponse struct {
	EmailData EmailData `json:"data"`
}

type PrincipalData struct {
	ID int64 `json:"id"`
}

type OrgResponse struct {
	OrgData OrgData `json:"org"`
}

type OrgData struct {
	Identifier string `json:"identifier"`
}

type ProjectResponse struct {
	ProjectData ProjectData `json:"project"`
}

type ProjectData struct {
	Identifier string `json:"identifier"`
}

type RepoData struct {
	Identifier string `json:"identifier"`
}

type PRData struct {
	Number           int    `json:"number"`
	Title            string `json:"title"`
	State            string `json:"state"`
	SourceSHA        string `json:"source_sha"`
	MergeCheckStatus string `json:"merge_check_status"`
}

type PRDetailsData struct {
	Type             string `json:"type"`
	Title            string `json:"title"`
	State            string `json:"state"`
	MergeCheckStatus string `json:"merge_check_status"`
}

type PRActivity struct {
	PRActivityDecision PRActivityDecision `json:"payload"`
	PRActivityAuthor   PRActivityAuthor   `json:"author"`
	Type               string             `json:"type"`
}

type PRActivityDecision struct {
	Decision *string `json:"decision"`
}
type PRActivityAuthor struct {
	DisplayName string `json:"display_name"`
}

type PRMergeRequest struct {
	BypassRules bool   `json:"bypass_rules"`
	DryRun      bool   `json:"dry_run"`
	SourceSHA   string `json:"source_sha"`
}

type PRMergeResponse struct {
	RuleViolations []PRRuleViolation `json:"rule_violations"`
}

type PRRuleViolation struct {
	Bypassable bool `json:"bypassable"`
}
