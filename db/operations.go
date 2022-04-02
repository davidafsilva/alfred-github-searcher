package db

import (
	"errors"
	"time"

	"github.com/davidafsilva/alfred-github-top-repositories/db/github"
	"github.com/davidafsilva/alfred-github-top-repositories/db/persistence"
	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/keychain"
)

const (
	githubTokenKey            = "agr_github_token"
	prsRefreshIntervalKey     = "agr_prs_refresh_interval"
	defaultPrsRefreshInterval = 30 * time.Minute
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
	return load(d.repositoriesFile)
}

func (d *Database) GetAllCreatedPRs() ([]persistence.PullRequest, error) {
	return d.loadAndRefreshPrsData(d.createdPrsFile, d.RefreshCreatedPRs)
}

func (d *Database) GetAllRequestedReviewPRs() ([]persistence.PullRequest, error) {
	return d.loadAndRefreshPrsData(d.requestedReviewPrsFile, d.RefreshRequestedReviewPRs)
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

func load[T any](jf *persistence.JsonFile[T]) ([]T, error) {
	// load data into the database
	if err := jf.Load(); err != nil {
		return nil, err
	}
	return jf.Data, nil
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

func (d *Database) loadAndRefreshPrsData(
	file *persistence.JsonFile[persistence.PullRequest],
	refreshFn func() ([]persistence.PullRequest, error),
) ([]persistence.PullRequest, error) {
	data, err := load(file)
	if err != nil {
		return nil, err
	}

	refreshInterval := d.wf.Config.GetDuration(prsRefreshIntervalKey, defaultPrsRefreshInterval)
	cacheExpirationDate := file.LastUpdated.Add(refreshInterval)
	if cacheExpirationDate.After(time.Now().UTC()) {
		// data needs to be refreshed
		return refreshFn()
	}

	return data, nil
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
