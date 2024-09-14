package main

import (
	"context"
	"fmt"
	"github.com/dhruv1397/pr-monitor/clientbuilder"
	"github.com/dhruv1397/pr-monitor/types"
	"log"
	"slices"
	"strconv"
	"strings"
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

func main() {
	githubPAT := ""
	//state := "closed"

	ctx := context.Background()

	client, err := clientbuilder.GetGithubPRClient(ctx, githubPAT)
	if err != nil {
		log.Fatal(err)
	}

	prs, err := client.GetPullRequests(ctx, nil, ConvertToPrintable)
	if err != nil {
		log.Fatal(err)
	}

	PrintPullRequests(prs)

}

func ConvertToPrintable(pr *types.PullRequest) *types.PrintablePullRequest {
	wrappedTitle := WrapText(pr.Title, colWidthTitle)
	wrappedPRNumber := WrapText(strconv.Itoa(pr.Number), colWidthPRNumber)
	wrappedSCMType := WrapText(pr.SCMProviderType, colWidthSCMType)
	wrappedState := WrapText(pr.State, colWidthState)
	wrappedMergeable := WrapText(pr.Mergeable, colWidthMergeable)
	wrappedApproved := WrapTextSlice(pr.Approved, colWidthApproved)
	wrappedCommented := WrapTextSlice(pr.Commented, colWidthCommented)
	wrappedRequestedChanges := WrapTextSlice(pr.RequestedChanges, colWidthRequestedChanges)
	wrappedURL := WrapText(pr.URL, colWidthURL)

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

func PrintPullRequests(prs []*types.PrintablePullRequest) {
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

func WrapText(text string, maxWidth int) []string {
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

func WrapTextSlice(text []string, maxWidth int) []string {
	return WrapText(strings.Join(text, ", "), maxWidth)
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
