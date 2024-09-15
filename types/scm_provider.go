package types

type SCMConfig struct {
	Providers []*SCMProvider `yaml:"scm_providers"`
}

type SCMProvider struct {
	Type    string  `yaml:"type"`
	Name    string  `yaml:"name"`
	Host    string  `yaml:"host"`
	User    *User   `yaml:"user"`
	Repos   []*Repo `yaml:"repos"`
	Updated int64   `yaml:"updated"`
	Created int64   `yaml:"created"`
}

type User struct {
	Name        string `yaml:"name"`
	PAT         string `yaml:"pat"`
	PrincipalID int64  `yaml:"principal_id"`
	Email       string `yaml:"email"`
}

type Repo struct {
	AccountIdentifier string `yaml:"account_identifier"`
	OrgIdentifier     string `yaml:"org_identifier"`
	ProjectIdentifier string `yaml:"project_identifier"`
	RepoIdentifier    string `yaml:"repo_identifier"`
}
