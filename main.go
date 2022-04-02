package main

import (
	"fmt"
	"os"

	"github.com/davidafsilva/alfred-github-top-repositories/action"
	"github.com/davidafsilva/alfred-github-top-repositories/action/repository"
	aw "github.com/deanishe/awgo"
	awu "github.com/deanishe/awgo/update"
)

var (
	githubRepository = "davidafsilva/alfred-github-repositories"
	githubIssues     = fmt.Sprintf("%s/issues", githubRepository)
	wf               *aw.Workflow
)

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
	case "search":
		err = repository.Search(wf, wf.Args()[1])
	case "sync":
		err = repository.Sync(wf)
	case "update":
		err = action.Update(wf)
	default:
		wf.Feedback.NewItem(fmt.Sprintf("Unknown action: %s", name)).
			Icon(aw.IconError)
	}

	if err != nil {
		wf.Feedback.NewItem(err.Error()).Icon(aw.IconError)
	}

	wf.SendFeedback()
}

func main() {
	if len(wf.Args()) < 1 {
		os.Exit(1)
	}

	wf.Run(run)
}
