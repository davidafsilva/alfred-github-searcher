package pullrequests

import (
	"fmt"
	"log"
	"sort"

	"github.com/davidafsilva/alfred-github-top-repositories/db"
	"github.com/davidafsilva/alfred-github-top-repositories/db/persistence"
	aw "github.com/deanishe/awgo"
)

var (
	iconPrDraft = &aw.Icon{
		Value: "pr-draft.png",
		Type:  aw.IconTypeImage,
	}
	iconPrReady = &aw.Icon{
		Value: "pr-ready.png",
		Type:  aw.IconTypeImage,
	}
)

func Search(wf *aw.Workflow, prFilter string) error {
	log.Println(fmt.Sprintf("executing pr search action with filter: %s", prFilter))

	database := db.New(wf)

	// get created pull requests
	createdPrs, err := database.GetAllCreatedPRs()
	if err != nil {
		return err
	}

	// get pull requests pending review
	pendingReviewPrs, err := database.GetAllPRsPendingReview()
	if err != nil {
		return err
	}

	// merge Prs
	allPrs := mergePullRequests(createdPrs, pendingReviewPrs)

	// add one item per pr
	for _, pr := range allPrs {
		item := wf.Feedback.NewItem(pr.Title).
			Subtitle(fmt.Sprintf("#%d opened by %s at %s", pr.Number, pr.Author, pr.Author)).
			Arg(pr.Url).
			Valid(true)
		if pr.IsDraft {
			item.Icon(iconPrDraft)
		} else {
			item.Icon(iconPrReady)
		}
	}

	warnEmptySubtitle := "Hint: "
	if prFilter != "" {
		wf.Filter(prFilter)
		warnEmptySubtitle += "Try a different search pattern or refresh " +
			"the pull requests with 'ghr'"
	} else {
		warnEmptySubtitle += "Try to refresh the pull requests with 'ghr'"
	}

	// fallback item when there are no pull requests
	if wf.IsEmpty() {
		wf.Feedback.NewItem("No pull requests found").
			Subtitle(warnEmptySubtitle).
			Icon(aw.IconWorkflow)
	}

	return nil
}

func mergePullRequests(created, toReview []persistence.PullRequest) []persistence.PullRequest {
	// merge them
	allPrs := make([]persistence.PullRequest, len(created)+len(toReview))
	copy(allPrs, created)
	copy(allPrs[len(created)-1:], toReview)

	// sort by creation date
	sort.Slice(allPrs, func(i, j int) bool {
		return allPrs[i].CreationDate.Before(allPrs[j].CreationDate)
	})

	return allPrs
}
