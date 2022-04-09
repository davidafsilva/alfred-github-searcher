package pullrequests

import (
	"fmt"
	"log"
	"sort"

	"github.com/davidafsilva/alfred-github-top-repositories/db"
	"github.com/davidafsilva/alfred-github-top-repositories/db/persistence"
	"github.com/davidafsilva/alfred-github-top-repositories/theme"
	aw "github.com/deanishe/awgo"
)

func Search(wf *aw.Workflow, prFilter string) error {
	return doSearch(wf, prFilter, func(database *db.Database) ([]persistence.PullRequest, error) {
		createdPrs, err := database.GetAllCreatedPRs()
		if err != nil {
			return nil, err
		}

		pendingReviewPrs, err := database.GetAllPRsPendingReview()
		if err != nil {
			return nil, err
		}

		return mergePullRequests(createdPrs, pendingReviewPrs), nil
	})
}

func SearchCreated(wf *aw.Workflow, prFilter string) error {
	return doSearch(wf, prFilter, func(database *db.Database) ([]persistence.PullRequest, error) {
		return database.GetAllCreatedPRs()
	})
}

func SearchPendingReview(wf *aw.Workflow, prFilter string) error {
	return doSearch(wf, prFilter, func(database *db.Database) ([]persistence.PullRequest, error) {
		return database.GetAllPRsPendingReview()
	})
}

func doSearch(
	wf *aw.Workflow,
	prFilter string,
	prsLoader func(database *db.Database) ([]persistence.PullRequest, error),
) error {
	log.Println(fmt.Sprintf("executing pr search action with filter: %s", prFilter))

	database := db.New(wf)

	// get created pull requests
	allPrs, err := prsLoader(database)
	if err != nil {
		return err
	}

	// add one item per pr
	icons := theme.New(wf).Icons
	for _, pr := range allPrs {
		item := wf.Feedback.NewItem(pr.Title).
			Subtitle(fmt.Sprintf("#%d opened by %s at %s", pr.Number, pr.Author, pr.RepositoryName)).
			Arg(pr.Url).
			Valid(true)
		if pr.IsDraft {
			item.Icon(icons.DraftPullRequest)
		} else {
			item.Icon(icons.PullRequest)
		}
	}

	warnEmptySubtitle := "Hint: "
	if prFilter != "" {
		wf.Filter(prFilter)
		warnEmptySubtitle += "Try a different search pattern or refresh " +
			"the pull requests with 'prsync'"
	} else {
		warnEmptySubtitle += "Try to refresh the pull requests with 'prsync'"
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
	copy(allPrs[len(created):], toReview)

	// sort by creation date
	sort.Slice(allPrs, func(i, j int) bool {
		return allPrs[i].CreationDate.After(allPrs[j].CreationDate)
	})

	return allPrs
}
