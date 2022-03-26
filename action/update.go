package action

import (
	"log"

	aw "github.com/deanishe/awgo"
)

func Update(wf *aw.Workflow) error {
	// check for newer releases
	if wf.UpdateCheckDue() {
		log.Println("checking for newer releases..")
		if err := wf.CheckForUpdate(); err != nil {
			log.Printf("error while checking for updates: %s\n", err.Error())
			return err
		}
	}

	// check if an update is available
	if !wf.UpdateAvailable() {
		wf.Feedback.NewItem("You're already on the latest release :)").
			Subtitle(wf.Version()).
			Icon(aw.IconWorkflow)
		return nil
	}

	// run the update
	if err := wf.InstallUpdate(); err != nil {
		log.Printf("error while updating to the latest version: %s\n", err.Error())
		return err
	}

	wf.Feedback.NewItem("Please accept and install the prompted workflow version").
		Icon(aw.IconWorkflow)
	return nil
}
