package db

import (
	"errors"
	"log"
	"time"

	"github.com/davidafsilva/alfred-github-top-repositories/db/github"
	"github.com/davidafsilva/alfred-github-top-repositories/db/persistence"
	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/keychain"
)

const (
	githubTokenKey                     = "ags_github_token"
	prsRefreshIntervalKey              = "ags_prs_refresh_interval"
	defaultPrsRefreshInterval          = 30 * time.Minute
	repositoriesRefreshIntervalKey     = "ags_repositories_refresh_interval"
	defaultRepositoriesRefreshInterval = 24 * 5 * time.Hour
)

type Database struct {
	wf                     *aw.Workflow
	repositoriesFile       *persistence.JsonFile[persistence.Repository]
	createdPrsFile         *persistence.JsonFile[persistence.PullRequest]
	requestedReviewPrsFile *persistence.JsonFile[persistence.PullRequest]
}

func New(wf *aw.Workflow) *Database {
	return &Database{
		wf:                     wf,
		repositoriesFile:       persistence.NewRepositoryJsonFile(wf),
		createdPrsFile:         persistence.NewCreatedPullRequestsJsonFile(wf),
		requestedReviewPrsFile: persistence.NewReviewRequestedPullRequestsJsonFile(wf),
	}
}

func (d *Database) GetAllRepositories() ([]persistence.Repository, error) {
	return loadAndRefreshData(
		d.wf,
		d.repositoriesFile,
		repositoriesRefreshIntervalKey,
		defaultRepositoriesRefreshInterval,
		d.RefreshRepositories,
	)
}

func (d *Database) GetAllCreatedPRs() ([]persistence.PullRequest, error) {
	return loadAndRefreshData(
		d.wf,
		d.createdPrsFile,
		prsRefreshIntervalKey,
		defaultPrsRefreshInterval,
		d.RefreshCreatedPRs,
	)
}

func (d *Database) GetAllPRsPendingReview() ([]persistence.PullRequest, error) {
	return loadAndRefreshData(
		d.wf,
		d.requestedReviewPrsFile,
		prsRefreshIntervalKey,
		defaultPrsRefreshInterval,
		d.RefreshRequestedReviewPRs,
	)
}

func (d *Database) RefreshRepositories() ([]persistence.Repository, error) {
	return refresh[github.Repository, persistence.Repository](
		d.wf,
		d.repositoriesFile,
		func(client *github.Client) ([]github.Repository, error) {
			return client.GetRepositories()
		},
		d.mapRepositoryData,
	)
}

func (d *Database) RefreshRequestedReviewPRs() ([]persistence.PullRequest, error) {
	return refresh[github.PullRequest, persistence.PullRequest](
		d.wf,
		d.requestedReviewPrsFile,
		func(client *github.Client) ([]github.PullRequest, error) {
			return client.GetReviewRequestedPullRequests()
		},
		mapPullRequestData,
	)
}

func (d *Database) RefreshCreatedPRs() ([]persistence.PullRequest, error) {
	return refresh[github.PullRequest, persistence.PullRequest](
		d.wf,
		d.createdPrsFile,
		func(client *github.Client) ([]github.PullRequest, error) {
			return client.GetCreatedPullRequests()
		},
		mapPullRequestData,
	)
}

func loadAndRefreshData[T any](
	wf *aw.Workflow,
	file *persistence.JsonFile[T],
	refreshConfigKey string,
	defaultRefreshInterval time.Duration,
	refreshFn func() ([]T, error),
) ([]T, error) {
	// load data into the database
	if err := file.Load(); err != nil {
		return nil, err
	}

	refreshInterval := wf.Config.GetDuration(refreshConfigKey, defaultRefreshInterval)
	cacheExpirationDate := file.LastUpdated.Add(refreshInterval)
	log.Printf("local data will expire at %v (%v TTL)\n", cacheExpirationDate, refreshInterval)
	if cacheExpirationDate.Before(time.Now().UTC()) {
		log.Println("data is stale, synchronizing..")
		// data needs to be refreshed
		return refreshFn()
	}

	return file.Data, nil
}

func refresh[I any, O any](
	wf *aw.Workflow,
	file *persistence.JsonFile[O],
	loadFn func(client *github.Client) ([]I, error),
	mappingFn func(entry I) O,
) ([]O, error) {
	// load the remote data
	token, err := resolveGitHubToken(wf)
	if err != nil {
		return nil, err
	}

	ghData, err := loadFn(github.NewClient(token))
	if err != nil {
		return nil, err
	}

	// map the repositories and download images
	data := make([]O, len(ghData))
	for i, ghEntry := range ghData {
		data[i] = mappingFn(ghEntry)
	}

	// update the JSON file
	if err = file.Save(data); err != nil {
		return nil, err
	}

	return file.Data, nil
}

func (d *Database) mapRepositoryData(entry github.Repository) persistence.Repository {
	var ownerImagePath string
	if len(entry.OpenGraphImageUrl) > 0 {
		ownerImagePath = downloadImage(d.wf, entry.OpenGraphImageUrl)
	}

	return persistence.Repository{
		Url:            entry.Url,
		Name:           entry.NameWithOwner,
		Description:    entry.Description,
		OwnerImagePath: ownerImagePath,
	}
}

func mapPullRequestData(entry github.PullRequest) persistence.PullRequest {
	return persistence.PullRequest{
		RepositoryName: entry.Repository.NameWithOwner,
		Url:            entry.Url,
		Title:          entry.Title,
		Number:         entry.Number,
		IsDraft:        entry.IsDraft,
		CreationDate:   entry.CreatedAt,
		Author:         entry.Author.LoginUser,
	}
}

func resolveGitHubToken(wf *aw.Workflow) (string, error) {
	// try env
	token := wf.Config.GetString(githubTokenKey)
	if len(token) > 0 {
		return token, nil
	}

	// try keychain
	token, err := wf.Keychain.Get(githubTokenKey)
	if err != nil && errors.Is(err, keychain.ErrNotFound) {
		err = errors.New("GitHub token was not found")
	}

	return token, err
}
