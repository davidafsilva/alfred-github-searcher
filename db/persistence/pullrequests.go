package persistence

import (
	aw "github.com/deanishe/awgo"
)

type PullRequest struct {
	RepositoryName string `json:"repositoryName"`
	Url            string `json:"url"`
	Title          string `json:"title"`
	Number         string `json:"number"`
	Author         string `json:"author"`
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
