package repository

import (
	"fmt"
	"log"

	"github.com/davidafsilva/alfred-github-top-repositories/db"
	aw "github.com/deanishe/awgo"
)

func Search(wf *aw.Workflow, repoFilter string) error {
	log.Println(fmt.Sprintf("executing search action with filter: %s", repoFilter))

	// get repositories
	database := db.New(wf)
	repositories, err := database.GetAllRepositories()
	if err != nil {
		return err
	}

	// add one item per repository
	showOwnerImages := wf.Config.GetBool("agr_show_owner_image", true)
	for _, r := range repositories {
		item := wf.Feedback.NewItem(r.Name).
			Subtitle(r.Description).
			Arg(r.Url).
			UID(r.Name).
			Valid(true)
		if showOwnerImages {
			item.Icon(&aw.Icon{
				Value: r.OwnerImagePath,
				Type:  aw.IconTypeImage,
			})
		} else {
			item.Icon(aw.IconWorkflow)
		}
	}

	warnEmptySubtitle := "Hint: "
	if repoFilter != "" {
		wf.Filter(repoFilter)
		warnEmptySubtitle += "Try a different search pattern or synchronize " +
			"the repositories with 'ghs'"
	} else {
		warnEmptySubtitle += "Try to synchronize the repositories with 'ghs'"
	}

	// fallback item when there are no persistence
	if wf.IsEmpty() {
		wf.Feedback.NewItem("No repositories found").
			Subtitle(warnEmptySubtitle).
			Icon(aw.IconWorkflow)
	}

	return nil
}
