package cli

import (
	"os"
)

const (
	CommandAdd     = "add"
	CommandRemove  = "remove"
	CommandRefresh = "refresh"
	CommandList    = "list"

	CommandAddHelpText     = "Add a new SCM provider."
	CommandRemoveHelpText  = "Remove a new SCM provider."
	CommandRefreshHelpText = "Refresh fetched values eg repos, user-name, etc for all the SCM providers."
	CommandListHelpText    = "List pull requests or SCM providers."

	SubcommandProvider  = "provider"
	SubcommandProviders = "providers"
	SubcommandPRs       = "prs"

	SubcommandAddProviderHelpText      = "Add an SCM provider."
	SubcommandRemoveProviderHelpText   = "Remove an SCM provider."
	SubcommandListProvidersHelpText    = "List SCM providers."
	SubcommandRefreshProvidersHelpText = "Refresh all the SCM providers."
	SubcommandPRsHelpText              = "List pull requests."

	ArgName         = "name"
	ArgNameHelpText = "Name of the SCM provider."

	FlagName   = "name"
	FlagType   = "type"
	FlagHost   = "host"
	FlagState  = "state"
	FlagOutput = "output"

	FlagNameShort   = 'n'
	FlagTypeShort   = 't'
	FlagHostShort   = 'h'
	FlagStateShort  = 's'
	FlagOutputShort = 'o'

	FlagNameHelpText   = "Name of the SCM provider."
	FlagTypeHelpText   = "Type of the SCM provider:- [github/harness]."
	FlagHostHelpText   = "Host URL of the SCM provider, eg https://github.com, https://app.harness.io."
	FlagStateHelpText  = "State of the pull request:- [open/merged/closed/all]"
	FlagOutputHelpText = "Output format:- [table/json/yaml]"
)

func GetArguments() []string {
	args := os.Args[1:]
	return args
}
