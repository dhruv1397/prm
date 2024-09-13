package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v64/github"
	"golang.org/x/oauth2"
)

func main() {
	// Replace with your GitHub personal access token (PAT)
	githubPAT := ""

	// Authenticate using the token
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubPAT},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// user, _, err := client.Users.Get(ctx, "")
	// if err != nil {
	// 	log.Fatalf("Error fetching authenticated user: %v", err)
	// }

	// // Print the user's login (username)
	// fmt.Printf("Authenticated GitHub Username: %s\n", user.GetLogin())

	// Construct the search query
	// query := "state:open author:dhruv-harness type:pr"

	// // Options for the search query
	// opts := &github.SearchOptions{
	// 	ListOptions: github.ListOptions{PerPage: 100}, // Adjust the number of results per page
	// }

	// // Execute the search query for issues (which includes pull requests)
	// result, _, err := client.Search.Issues(ctx, query, opts)
	// if err != nil {
	// 	log.Fatalf("Error executing search query: %v", err)
	// }

	// // Iterate through the search results and print relevant information
	// for _, issue := range result.Issues {
	// 	fmt.Printf("PR #%d: %s\nURL: %s\nState: %s\nAuthor: %s\n\n",
	// 		*issue.Number, *issue.Title, *issue.HTMLURL, *issue.State, *issue.User.Login)
	// }

		// Fetch the reviews for the PR
		reviews, _, err := client.PullRequests.ListReviews(ctx, "harness", "harness-core", 55964, nil)
		if err != nil {
			log.Fatalf("Error fetching reviews: %v", err)
		}
	
		// Print review statuses
		for _, review := range reviews {
			fmt.Printf("Review by %s: %s (%s)\n", review.GetUser().GetLogin(), review.GetState(), review.GetHTMLURL())
		}

		pr, _, err := client.PullRequests.Get(ctx,"harness", "harness-core", 55964)
		if err != nil {
			log.Fatalf("Error fetching PR details: %v", err)
		}
	
		// Check mergeability status
		if pr.Mergeable != nil && *pr.Mergeable {
			fmt.Println("The PR is mergeable.")
		} else {
			fmt.Println("The PR is not mergeable.")
		}
}
