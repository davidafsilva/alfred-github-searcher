package github

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shurcooL/githubv4"
)

type PullRequest struct {
	Repository Repository
	Url        string
	Title      string
	Number     int
	IsDraft    bool
	CreatedAt  time.Time
	Author     struct {
		LoginUser string `graphql:"login"`
	}
}

type pullRequestsSearchQuery struct {
	Search struct {
		PageInfo struct {
			HasNextPage bool
			EndCursor   githubv4.String
		}
		PullRequests []struct {
			PullRequest `graphql:"... on PullRequest"`
		} `graphql:"nodes"`
	} `graphql:"search(type: ISSUE, first: 100, after: $pageCursor, query: $searchQuery)"`
}

func (c *Client) GetReviewRequestedPullRequests() ([]PullRequest, error) {
	return c.getPullRequests("review-requested:@me")
}

func (c *Client) GetCreatedPullRequests() ([]PullRequest, error) {
	return c.getPullRequests("author:@me")
}

func (c *Client) getPullRequests(query string) ([]PullRequest, error) {
	var pullRequests []PullRequest
	var pageCursor *githubv4.String
	pageNumber := 1
	for {
		sq := fmt.Sprintf("is:open is:pr archived:false %s", query)
		variables := map[string]interface{}{
			"pageCursor":  pageCursor,
			"searchQuery": githubv4.String(sq),
		}
		var q pullRequestsSearchQuery
		log.Printf("requesting pull requests from GitHub (page %d)..\n", pageNumber)
		err := c.client.Query(context.Background(), &q, variables)
		if err != nil {
			log.Printf("error while contacting GitHub: %s\n", err.Error())
			return nil, err
		}

		for _, pr := range q.Search.PullRequests {
			pullRequests = append(pullRequests, pr.PullRequest)
		}
		if !q.Search.PageInfo.HasNextPage {
			break
		}
		pageCursor = githubv4.NewString(q.Search.PageInfo.EndCursor)
		pageNumber++
	}

	log.Printf("collected %d pull requests from GitHub\n", len(pullRequests))
	return pullRequests, nil
}
