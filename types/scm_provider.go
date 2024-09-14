package types

type SCMProvider struct {
	Identifier string
	Type       string // GITHUB, HARNESS
	Alias      string
	Host       string
	Token      string
}
