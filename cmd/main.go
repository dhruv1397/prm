package main

import (
	"context"
	"fmt"
	"github.com/dhruv1397/pr-monitor/clientbuilder"
	"github.com/dhruv1397/pr-monitor/types"
	"log"
	"strconv"
	"strings"
)

func main() {
	githubPAT := ""

	ctx := context.Background()

	client, err := clientbuilder.GetGithubClient(ctx, githubPAT)
	if err != nil {
		log.Fatal(err)
	}

	prs, err := client.GetOpenPullRequests(ctx)
	if err != nil {
		log.Fatal(err)
	}

	PrintPullRequests(prs)

}

func PrintPullRequests(prs []*types.PullRequest) {
	const colWidthSerialNumber = 4
	const colWidthTitle = 34
	const colWidthPRNumber = 10
	const colWidthSCMType = 10
	const colWidthMergeable = 10
	const colWidthApproved = 17
	const colWidthCommented = 17
	const colWidthRequestedChanges = 17
	const colWidthURL = 34
	const spacingPattern = "| %-4s | %-34s | %-10s | %-10s | %-10s | %-17s | %-17s | %-17s | %-34s |\n"
	const separatorLength = 28 + colWidthSerialNumber + colWidthTitle + colWidthPRNumber +
		colWidthSCMType + colWidthMergeable + colWidthApproved +
		colWidthCommented + colWidthRequestedChanges + colWidthURL
	printSeparator(separatorLength)

	fmt.Printf(spacingPattern, "#", "Title", "PR Number", "SCM Type", "Mergeable", "Approved", "Commented",
		"Requested Changes", "URL")

	printSeparator(separatorLength)

	for index, pr := range prs {

		wrappedSerialNumber := wrapText(strconv.Itoa(index), colWidthSerialNumber)
		wrappedTitle := wrapText(pr.Title, colWidthTitle)
		wrappedPRNumber := wrapText(strconv.Itoa(pr.Number), colWidthPRNumber)
		wrappedSCMType := wrapText(pr.SCMProviderType, colWidthSCMType)
		wrappedMergeable := wrapText(strconv.FormatBool(pr.Mergeable), colWidthMergeable)
		wrappedApproved := wrapTextSlice(pr.Approved, colWidthApproved)
		wrappedCommented := wrapTextSlice(pr.Commented, colWidthCommented)
		wrappedRequestedChanges := wrapTextSlice(pr.ChangesRequested, colWidthRequestedChanges)
		wrappedURL := wrapText(pr.URL, colWidthURL)

		maxRows := max(
			len(wrappedSerialNumber),
			len(wrappedTitle),
			len(wrappedPRNumber),
			len(wrappedSCMType),
			len(wrappedMergeable),
			len(wrappedApproved),
			len(wrappedCommented),
			len(wrappedRequestedChanges),
			len(wrappedURL),
		)

		for i := 0; i < maxRows; i++ {
			fmt.Printf(spacingPattern,
				getListElement(wrappedSerialNumber, i),
				getListElement(wrappedTitle, i),
				getListElement(wrappedPRNumber, i),
				getListElement(wrappedSCMType, i),
				getListElement(wrappedMergeable, i),
				getListElement(wrappedApproved, i),
				getListElement(wrappedCommented, i),
				getListElement(wrappedRequestedChanges, i),
				getListElement(wrappedURL, i),
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
