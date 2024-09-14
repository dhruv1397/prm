package types

type SCMProviderGithub struct {
	User GithubUser `yaml:"user"`
	SCMProvider
}

type GithubUser struct {
	Name string `yaml:"name"`
}
