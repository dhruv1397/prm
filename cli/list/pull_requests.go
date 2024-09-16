package list

import (
	"context"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"github.com/dhruv1397/pr-monitor/cli"
	"github.com/dhruv1397/pr-monitor/clientbuilder"
	"github.com/dhruv1397/pr-monitor/prclient"
	"github.com/dhruv1397/pr-monitor/store"
	"github.com/dhruv1397/pr-monitor/types"
	"github.com/dhruv1397/pr-monitor/util"
	"sync"

	"slices"
	"strconv"
	"strings"

	"time"
)

const (
	colWidthSerialNumber     = 4
	colWidthTitle            = 34
	colWidthPRNumber         = 10
	colWidthSCMType          = 10
	colWidthState            = 10
	colWidthMergeable        = 10
	colWidthApproved         = 17
	colWidthCommented        = 17
	colWidthRequestedChanges = 17
	colWidthURL              = 34
	separatorLength          = 31 + colWidthSerialNumber + colWidthTitle + colWidthPRNumber +
		colWidthSCMType + colWidthState + colWidthMergeable + colWidthApproved +
		colWidthCommented + colWidthRequestedChanges + colWidthURL
	spacingPattern = "| %-4s | %-34s | %-10s | %-10s | %-10s | %-10s | %-17s | %-17s | %-17s | %-34s |\n"
)

type prsCommand struct {
	state        string
	providerType string
	providerName string
}

func (c *prsCommand) run(*kingpin.ParseContext) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	str := store.NewSCMProviderImpl()
	providers, err := str.List(c.providerType, c.providerName)
	if err != nil {
		return fmt.Errorf("failed to list providers: %w", err)
	}
	if len(providers) == 0 {
		fmt.Println("No providers found!")
		return nil
	}
	return c.helper(ctx, providers)
}

func (c *prsCommand) helper(ctx context.Context, providers []*types.SCMProvider) error {
	var allPRs = make([]*types.PrintablePullRequest, 0)
	var errs []error

	respCh := make(chan []*types.PrintablePullRequest)
	errCh := make(chan error)
	var wg sync.WaitGroup

	var prMutex sync.Mutex
	var errMutex sync.Mutex

	for _, provider := range providers {
		wg.Add(1)
		go func(provider *types.SCMProvider) {
			defer wg.Done()

			prClient, err := c.getPRClient(ctx, provider)
			if err != nil {
				errCh <- err
				return
			}

			prs, err := prClient.GetPullRequests(ctx, c.state, ConvertToPrintable)
			if err != nil {
				errCh <- err
				return
			}

			respCh <- prs
		}(provider)
	}

	go func() {
		wg.Wait()
		close(respCh)
		close(errCh)
	}()

	for respCh != nil || errCh != nil {
		select {
		case prs, ok := <-respCh:
			if !ok {
				respCh = nil
			} else {
				prMutex.Lock()
				allPRs = append(allPRs, prs...)
				prMutex.Unlock()
			}
		case err, ok := <-errCh:
			if !ok {
				errCh = nil
			} else {
				errMutex.Lock()
				errs = append(errs, err)
				errMutex.Unlock()
			}
		}
	}

	if len(allPRs) > 0 {
		printPullRequests(allPRs)
	} else {
		fmt.Println("No PRs found!")
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors encountered:\n%v", util.FormatErrors(errs))
	}

	return nil
}

func (c *prsCommand) getPRClient(ctx context.Context, provider *types.SCMProvider) (prclient.PRClient, error) {
	if provider.Type == "github" {
		return clientbuilder.GetGithubPRClient(ctx, provider.User)
	} else if provider.Type == "harness" {
		return clientbuilder.GetHarnessPRClient(provider.Host, provider.User, provider.Repos)
	} else {
		return nil, fmt.Errorf("unknown provider type: %s", provider.Type)
	}
}

func registerPRs(app *kingpin.CmdClause) {
	c := &prsCommand{}

	cmd := app.Command(cli.SubcommandPRs, cli.SubcommandPRsHelpText).Default().Action(c.run)

	cmd.Flag(cli.FlagState, cli.FlagStateHelpText).Default("open").StringVar(&c.state)

	cmd.Flag(cli.FlagType, cli.FlagTypeHelpText).StringVar(&c.providerType)

	cmd.Flag(cli.FlagName, cli.FlagNameHelpText).StringVar(&c.providerName)
}

func ConvertToPrintable(pr *types.PullRequest) *types.PrintablePullRequest {
	wrappedTitle := wrapText(pr.Title, colWidthTitle)
	wrappedPRNumber := wrapText(strconv.Itoa(pr.Number), colWidthPRNumber)
	wrappedSCMType := wrapText(pr.SCMProviderType, colWidthSCMType)
	wrappedState := wrapText(pr.State, colWidthState)
	wrappedMergeable := wrapText(pr.Mergeable, colWidthMergeable)
	wrappedApproved := wrapTextSlice(pr.Approved, colWidthApproved)
	wrappedCommented := wrapTextSlice(pr.Commented, colWidthCommented)
	wrappedRequestedChanges := wrapTextSlice(pr.RequestedChanges, colWidthRequestedChanges)
	wrappedURL := wrapText(pr.URL, colWidthURL)

	maxRows := max(
		len(wrappedTitle),
		len(wrappedPRNumber),
		len(wrappedSCMType),
		len(wrappedState),
		len(wrappedMergeable),
		len(wrappedApproved),
		len(wrappedCommented),
		len(wrappedRequestedChanges),
		len(wrappedURL),
	)
	return &types.PrintablePullRequest{
		NumberRaw:          pr.Number,
		SCMProviderTypeRaw: pr.SCMProviderType,
		Title:              wrappedTitle,
		Number:             wrappedPRNumber,
		SCMProviderType:    wrappedSCMType,
		State:              wrappedState,
		Mergeable:          wrappedMergeable,
		Approved:           wrappedApproved,
		Commented:          wrappedCommented,
		RequestedChanges:   wrappedRequestedChanges,
		URL:                wrappedURL,
		MaxRows:            maxRows,
	}
}

func printPullRequests(prs []*types.PrintablePullRequest) {
	printSeparator(separatorLength)
	fmt.Printf(spacingPattern, "#", "Title", "PR Number", "SCM Type", "State", "Mergeable", "Approved", "Commented",
		"Requested Changes", "URL")
	printSeparator(separatorLength)

	slices.SortFunc(prs, types.ComparePrintablePullRequest)

	for index, pr := range prs {
		for i := 0; i < pr.MaxRows; i++ {
			fmt.Printf(spacingPattern,
				getSrNumberElement(index, i),
				getListElement(pr.Title, i),
				getListElement(pr.Number, i),
				getListElement(pr.SCMProviderType, i),
				getListElement(pr.State, i),
				getListElement(pr.Mergeable, i),
				getListElement(pr.Approved, i),
				getListElement(pr.Commented, i),
				getListElement(pr.RequestedChanges, i),
				getListElement(pr.URL, i),
			)
		}
		printSeparator(separatorLength)
	}
}

func printSeparator(length int) {
	fmt.Println(strings.Repeat("-", length))
}

func wrapText(text string, maxWidth int) []string {
	var chunks []string

	for i := 0; i < len(text); i += maxWidth {
		end := i + maxWidth
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, text[i:end])
	}

	return chunks
}

func wrapTextSlice(text []string, maxWidth int) []string {
	return wrapText(strings.Join(text, ", "), maxWidth)
}

func max(values ...int) int {
	maxValue := values[0]
	for _, value := range values {
		if value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}

func getListElement(text []string, index int) string {
	if index >= len(text) {
		return ""
	}
	return text[index]
}

func getSrNumberElement(srNumber int, index int) string {
	if index == 0 {
		return strconv.Itoa(srNumber)
	}
	return ""
}