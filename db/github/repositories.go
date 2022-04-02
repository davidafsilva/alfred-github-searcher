package github

import (
	"context"
	"log"

	"github.com/shurcooL/githubv4"
)

type Repository struct {
	NameWithOwner     string
	Description       string
	Url               string
	OpenGraphImageUrl string
}
type repositoriesQuery struct {
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

func (c *Client) GetRepositories() ([]Repository, error) {
	var repositories []Repository
	var pageCursor *githubv4.String
	pageNumber := 1
	for {
		variables := map[string]interface{}{
			"pageCursor": pageCursor,
		}
		var q repositoriesQuery
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
