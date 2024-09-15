package list

import (
	"context"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"github.com/dhruv1397/pr-monitor/clientbuilder"
	"github.com/dhruv1397/pr-monitor/prclient"
	"github.com/dhruv1397/pr-monitor/store"
	"github.com/dhruv1397/pr-monitor/types"
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
		return fmt.Errorf("failed to list providers: %v", err)
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

	// Create mutexes to synchronize access to the allPRs and errs slices
	var prMutex sync.Mutex
	var errMutex sync.Mutex

	// Iterate over providers and fetch PRs in parallel
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

	// Close channels after all goroutines are done
	go func() {
		wg.Wait()
		close(respCh)
		close(errCh)
	}()

	// Collect PRs and errors
	for respCh != nil || errCh != nil {
		select {
		case prs, ok := <-respCh:
			if !ok {
				respCh = nil
			} else {
				// Add PRs to the result list safely using prMutex
				prMutex.Lock()
				allPRs = append(allPRs, prs...)
				prMutex.Unlock()
			}
		case err, ok := <-errCh:
			if !ok {
				errCh = nil
			} else {
				// Add error to the error list safely using errMutex
				errMutex.Lock()
				errs = append(errs, err)
				errMutex.Unlock()
			}
		}
	}

	// Print the collected PRs
	if len(allPRs) > 0 {
		printPullRequests(allPRs)
	} else {
		fmt.Println("No PRs found!")
	}

	// If there are errors, return the first one (you can customize error handling here)
	if len(errs) > 0 {
		return fmt.Errorf("errors occurred: %v", errs)
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

	cmd := app.Command("prs", "list pull requests").Default().Action(c.run)

	cmd.Flag("state", "state of the pull request").Default("open").StringVar(&c.state)

	cmd.Flag("type", "type of the SCM provider").StringVar(&c.providerType)

	cmd.Flag("provider", "name of the SCM provider").StringVar(&c.providerName)
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
