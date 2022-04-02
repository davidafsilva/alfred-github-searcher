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

	actionSearch     = "search"
	actionRefresh    = "refresh"
	actionUpdate     = "update"
	targetRepository = "repository"
	targetPr         = "pr"
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
		err = search(wf.Args()[1], strings.Join(wf.Args()[2:], ""))
	case actionRefresh:
		err = refresh(wf.Args()[1])
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

func search(target string, query string) error {
	var err error = nil
	switch target {
	case targetRepository:
		err = repository.Search(wf, query)
	case targetPr:
		err = pullrequests.Search(wf, query)
	default:
		err = errors.New(fmt.Sprintf("Unknown action target: %s %s", actionSearch, target))
	}
	return err
}

func refresh(target string) error {
	var err error = nil
	switch target {
	case targetRepository:
		err = repository.Refresh(wf)
	case targetPr:
		err = pullrequests.Refresh(wf)
	default:
		err = errors.New(fmt.Sprintf("Unknown action target: %s %s", actionRefresh, target))
	}
	return err
}

func main() {
	wf.Run(run)
}
