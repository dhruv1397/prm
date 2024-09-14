package types

type SCMProvider struct {
	Identifier string `yaml:"identifier"`
	Type       string `yaml:"type"` // github, harness
	Alias      string `yaml:"alias,omitempty"`
	Host       string `yaml:"host"`
	Token      string `yaml:"token"`
}
