package action

import (
	"fmt"
	"log"

	"github.com/davidafsilva/alfred-github-top-repositories/db"
	aw "github.com/deanishe/awgo"
)

func Sync(wf *aw.Workflow) error {
	log.Println("executing sync action")

	database := db.NewDatabase(wf)
	repositories, err := database.Refresh()
	if err != nil {
		return err
	}

	wf.Feedback.NewItem("Successfully synchronized local database").
		Subtitle(fmt.Sprintf("%d repositories found", len(repositories))).
		Icon(aw.IconWorkflow)
	return nil
}
