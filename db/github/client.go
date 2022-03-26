package github

import (
	"context"
	"log"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Repository struct {
	NameWithOwner string `graphql:"nameWithOwner"`
	Description   string `graphql:"description"`
	Url           string `graphql:"url"`
}
type query struct {
	Viewer struct {
		TopRepositories struct {
			PageInfo struct {
				HasNextPage bool
				EndCursor   githubv4.String
			}
			Nodes []Repository
		} `graphql:"topRepositories(first: 100, after: $pageCursor, orderBy: {field: UPDATED_AT, direction: DESC})"`
	}
}

type Client struct {
	client *githubv4.Client
}

func NewClient(token string) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), ts)
	return &Client{
		client: githubv4.NewClient(httpClient),
	}
}

func (c *Client) GetRepositories() ([]Repository, error) {
	var repositories []Repository
	var pageCursor *githubv4.String
	pageNumber := 1
	for {
		variables := map[string]interface{}{
			"pageCursor": pageCursor,
		}
		var q query
		log.Printf("requesting repositories from GitHub (page %d)..\n", pageNumber)
		err := c.client.Query(context.Background(), &q, variables)
		if err != nil {
			log.Printf("error while contacting GitHub: %s\n", err.Error())
			return nil, err
		}

		repositories = append(repositories, q.Viewer.TopRepositories.Nodes...)
		if !q.Viewer.TopRepositories.PageInfo.HasNextPage {
			break
		}
		pageCursor = githubv4.NewString(q.Viewer.TopRepositories.PageInfo.EndCursor)
		pageNumber++
	}

	log.Printf("collected %d repositores from GitHub\n", len(repositories))
	return repositories, nil
}
