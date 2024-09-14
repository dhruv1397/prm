package types

type SCMProviderHarness struct {
	User  HarnessUser   `yaml:"user"`
	Repos []HarnessRepo `yaml:"repos"`
	SCMProvider
}

type HarnessRepo struct {
	AccountIdentifier string `yaml:"account_identifier"`
	OrgIdentifier     string `yaml:"org_identifier"`
	ProjectIdentifier string `yaml:"project_identifier"`
	RepoIdentifier    string `yaml:"repo_identifier"`
}

type HarnessUser struct {
	PrincipalID int64  `yaml:"principal_id"`
	Email       string `yaml:"email"`
	Token       string `yaml:"pat"`
}
