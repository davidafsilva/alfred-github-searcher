package persistence

import (
	"fmt"
	"log"
	"time"

	aw "github.com/deanishe/awgo"
)

type JsonFile[T any] struct {
	dataType    string
	wf          *aw.Workflow
	LastUpdated time.Time `json:"lastUpdated"`
	Data        []T       `json:"data"`
}

func (jf *JsonFile[T]) Load() error {
	log.Printf("loading local %s database from %s..\n", jf.dataType, jf.filename())

	// check if the data was already loaded
	if !jf.LastUpdated.IsZero() {
		log.Printf("%s are cached, skipped file read\n", jf.dataType)
		return nil
	}

	if err := jf.loadFromFile(); err != nil {
		log.Printf("error loading file: %s\n", err.Error())
		return err
	}

	log.Printf("local %s database loaded!\n", jf.dataType)
	return nil
}

func (jf *JsonFile[T]) Save(data []T) error {
	log.Printf("saving local %s database to %s..\n", jf.dataType, jf.filename())

	jf.Data = data
	jf.LastUpdated = time.Now().UTC()
	if err := jf.saveToFile(); err != nil {
		log.Printf("error saving file: %s\n", err.Error())
	}

	log.Printf("local %s database saved!\n", jf.dataType)
	return nil
}

func (jf *JsonFile[T]) loadFromFile() error {
	if !jf.wf.Data.Exists(jf.filename()) {
		return nil
	}

	return jf.wf.Data.LoadJSON(jf.filename(), jf)
}

func (jf *JsonFile[T]) saveToFile() error {
	return jf.wf.Data.StoreJSON(jf.filename(), jf)
}

func (jf *JsonFile[T]) filename() string {
	return fmt.Sprintf("%s.json", jf.dataType)
}
