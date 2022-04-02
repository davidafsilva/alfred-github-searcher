package persistence

import (
	"time"

	aw "github.com/deanishe/awgo"
)

type PullRequest struct {
	RepositoryName string    `json:"repositoryName"`
	Url            string    `json:"url"`
	Title          string    `json:"title"`
	Number         int       `json:"number"`
	Author         string    `json:"author"`
	IsDraft        bool      `json:"isDraft"`
	CreationDate   time.Time `json:"creationDate"`
}

func NewCreatedPullRequestsJsonFile(wf *aw.Workflow) *JsonFile[PullRequest] {
	return &JsonFile[PullRequest]{
		dataType: "pull-requests-created",
		wf:       wf,
	}
}

func NewReviewRequestedPullRequestsJsonFile(wf *aw.Workflow) *JsonFile[PullRequest] {
	return &JsonFile[PullRequest]{
		dataType: "pull-requests-review-requests",
		wf:       wf,
	}
}
