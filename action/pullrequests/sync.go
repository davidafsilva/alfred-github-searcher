package pullrequests

import (
	"fmt"
	"log"

	"github.com/davidafsilva/alfred-github-top-repositories/db"
	aw "github.com/deanishe/awgo"
)

func Sync(wf *aw.Workflow) error {
	log.Println("executing pull requests synchronization action..")

	database := db.New(wf)
	prs, err := database.RefreshCreatedPRs()
	if err != nil {
		return err
	}
	total := len(prs)

	prs, err = database.RefreshRequestedReviewPRs()
	if err != nil {
		return err
	}
	total += len(prs)

	wf.Feedback.NewItem("Successfully synchronized local database").
		Subtitle(fmt.Sprintf("%d pull requests found", total)).
		Icon(aw.IconWorkflow)
	return nil
}
