package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/davidafsilva/alfred-github-top-repositories/action"
	"github.com/davidafsilva/alfred-github-top-repositories/action/pullrequests"
	"github.com/davidafsilva/alfred-github-top-repositories/action/repository"
	aw "github.com/deanishe/awgo"
	awu "github.com/deanishe/awgo/update"
)

const (
	githubRepository = "davidafsilva/alfred-github-repositories"
	githubIssues     = "davidafsilva/alfred-github-repositories/issues"

	actionSearch = "search"
	actionSync   = "sync"
	actionUpdate = "update"

	targetRepository          = "repository"
	targetPr                  = "pr"
	targetPrTypeCreated       = "created"
	targetPrTypePendingReview = "pending-review"
)

var wf *aw.Workflow

func init() {
	wf = aw.New(
		awu.GitHub(githubRepository),
		aw.HelpURL(githubIssues),
		aw.MaxResults(10),
	)
}

func run() {
	name := wf.Args()[0]
	var err error = nil
	switch name {
	case actionSearch:
		err = search(wf.Args()[1], wf.Args()[2:])
	case actionSync:
		err = sync(wf.Args()[1])
	case actionUpdate:
		err = action.Update(wf)
	default:
		err = errors.New(fmt.Sprintf("Unknown action: %s", name))
	}

	if err != nil {
		wf.Feedback.NewItem(err.Error()).Icon(aw.IconError)
	}

	wf.SendFeedback()
}

func search(target string, args []string) error {
	var err error = nil
	switch target {
	case targetRepository:
		err = repository.Search(wf, strings.Join(args, ""))
	case targetPr:
		t := args[0]
		q := strings.Join(args[1:], "")
		switch t {
		case targetPrTypeCreated:
			err = pullrequests.SearchCreated(wf, q)
		case targetPrTypePendingReview:
			err = pullrequests.SearchPendingReview(wf, q)
		default:
			err = pullrequests.Search(wf, q)
		}
	default:
		err = errors.New(fmt.Sprintf("Unknown action target: %s %s", actionSearch, target))
	}
	return err
}

func sync(target string) error {
	var err error = nil
	switch target {
	case targetRepository:
		err = repository.Sync(wf)
	case targetPr:
		err = pullrequests.Sync(wf)
	default:
		err = errors.New(fmt.Sprintf("Unknown action target: %s %s", actionSync, target))
	}
	return err
}

func main() {
	wf.Run(run)
}
