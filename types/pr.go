package types

type PullRequest struct {
	Number           int      // 2234
	Title            string   // Doing something new
	SCMProviderType  string   // Github, Harness
	URL              string   // https://github.com/abc/abc/pul/123
	Approved         []string // ["user 1", "user 2]
	Commented        []string // ["user 1", "user 2]
	ChangesRequested []string // ["user 1", "user 2]
	Mergeable        bool     // true
}
