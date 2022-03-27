package db

import (
	"errors"

	"github.com/davidafsilva/alfred-github-top-repositories/db/github"
	"github.com/davidafsilva/alfred-github-top-repositories/db/persistence"
	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/keychain"
)

type Database struct {
	wf *aw.Workflow
	jf *persistence.JsonFile
}

func NewDatabase(wf *aw.Workflow) *Database {
	return &Database{
		wf: wf,
		jf: persistence.NewJsonFile(wf),
	}
}

func (d *Database) GetAllRepositories() ([]persistence.Repository, error) {
	// load data into the database
	if err := d.jf.Load(); err != nil {
		return nil, err
	}
	return d.jf.Repositories, nil
}

func (d *Database) Refresh() ([]persistence.Repository, error) {
	// load the remote data
	token, err := resolveGitHubToken(d.wf)
	if err != nil {
		return nil, err
	}

	ghRepos, err := github.NewClient(token).GetRepositories()
	if err != nil {
		return nil, err
	}

	// map the repositories and download images
	repos := make([]persistence.Repository, len(ghRepos))
	for i, ghr := range ghRepos {
		var ownerImagePath string
		if len(ghr.OpenGraphImageUrl) > 0 {
			ownerImagePath = downloadImage(d.wf, ghr.OpenGraphImageUrl)
		}

		repos[i] = persistence.Repository{
			Url:            ghr.Url,
			Name:           ghr.NameWithOwner,
			Description:    ghr.Description,
			OwnerImagePath: ownerImagePath,
		}
	}

	// update the JSON file
	if err = d.jf.Save(repos); err != nil {
		return nil, err
	}

	return d.jf.Repositories, nil
}

func resolveGitHubToken(wf *aw.Workflow) (string, error) {
	key := "alfred_github_repos_token"

	// try env
	token := wf.Config.GetString(key)
	if len(token) > 0 {
		return token, nil
	}

	// try keychain
	token, err := wf.Keychain.Get(key)
	if err != nil && errors.Is(err, keychain.ErrNotFound) {
		err = errors.New("GitHub token was not found")
	}

	return token, err
}
