package repository

import (
	"fmt"
	"log"

	"github.com/davidafsilva/alfred-github-top-repositories/db"
	"github.com/davidafsilva/alfred-github-top-repositories/theme"
	aw "github.com/deanishe/awgo"
)

const showOwnerImagesKey = "ags_show_owner_image"

func Search(wf *aw.Workflow, repoFilter string) error {
	log.Println(fmt.Sprintf("executing repository search action with filter: %s", repoFilter))

	// get repositories
	database := db.New(wf)
	repositories, err := database.GetAllRepositories()
	if err != nil {
		return err
	}

	// add one item per repository
	showOwnerImages := wf.Config.GetBool(showOwnerImagesKey, true)
	icons := theme.New(wf).Icons
	for _, r := range repositories {
		item := wf.Feedback.NewItem(r.Name).
			Subtitle(r.Description).
			Arg(r.Url).
			UID(r.Name).
			Valid(true)
		if !showOwnerImages {
			item.Icon(&aw.Icon{
				Value: r.OwnerImagePath,
				Type:  aw.IconTypeImage,
			})
		} else {
			item.Icon(icons.Repository)
		}
	}

	warnEmptySubtitle := "Hint: "
	if repoFilter != "" {
		wf.Filter(repoFilter)
		warnEmptySubtitle += "Try a different search pattern or sync " +
			"the repositories with 'reposync'"
	} else {
		warnEmptySubtitle += "Try to sync the repositories with 'reposync'"
	}

	// fallback item when there are no repositories
	if wf.IsEmpty() {
		wf.Feedback.NewItem("No repositories found").
			Subtitle(warnEmptySubtitle).
			Icon(aw.IconWorkflow)
	}

	return nil
}
