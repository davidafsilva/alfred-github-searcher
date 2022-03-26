package action

import (
	"fmt"
	"log"

	"github.com/davidafsilva/alfred-github-top-repositories/db"
	aw "github.com/deanishe/awgo"
)

func Search(wf *aw.Workflow, repoFilter string) error {
	log.Println(fmt.Sprintf("executing search action with filter: %s", repoFilter))

	// get repositories
	database := db.NewDatabase(wf)
	repositories, err := database.GetAllRepositories()
	if err != nil {
		return err
	}

	// add one item per repository
	for _, r := range repositories {
		wf.Feedback.NewItem(r.Name).
			Subtitle(r.Description).
			Icon(aw.IconWorkflow).
			Arg(r.Url).
			UID(r.Name).
			Valid(true)
	}

	warnEmptySubtitle := "Hint: Try to synchronize the repositories with 'ghs'"
	if repoFilter != "" {
		wf.Filter(repoFilter)
		warnEmptySubtitle = "Hint: Try a different search pattern or synchronize the repositories with 'ghs'"
	}

	// fallback item when there are no persistence
	if wf.IsEmpty() {
		wf.Feedback.NewItem("No repositories found").
			Subtitle(warnEmptySubtitle).
			Icon(aw.IconWorkflow)
	}

	return nil
}