package persistence

import (
	"log"
	"time"

	aw "github.com/deanishe/awgo"
)

type Repository struct {
	Url         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type JsonFile struct {
	wf           *aw.Workflow
	LastUpdated  time.Time    `json:"lastUpdated"`
	Repositories []Repository `json:"repositories"`
}

const filename = "database.json"

func NewJsonFile(wf *aw.Workflow) *JsonFile {
	return &JsonFile{wf: wf}
}

func (jf *JsonFile) Load() error {
	log.Printf("loading local database from %s..\n", filename)

	// check if the data was already loaded
	if !jf.LastUpdated.IsZero() {
		log.Println("repositories are cached, skipped file read")
		return nil
	}

	if err := jf.loadFromFile(); err != nil {
		log.Printf("error loading file: %s\n", err.Error())
		return err
	}

	log.Println("local database loaded!")
	return nil
}

func (jf *JsonFile) Save(repos []Repository) error {
	log.Printf("saving local database to %s..\n", filename)

	jf.Repositories = repos
	jf.LastUpdated = time.Now().UTC()
	if err := jf.saveToFile(); err != nil {
		log.Printf("error saving file: %s\n", err.Error())
	}

	log.Println("local database saved!")
	return nil
}

func (jf *JsonFile) loadFromFile() error {
	if !jf.wf.Data.Exists(filename) {
		return nil
	}

	return jf.wf.Data.LoadJSON(filename, jf)
}

func (jf *JsonFile) saveToFile() error {
	return jf.wf.Data.StoreJSON(filename, jf)
}
