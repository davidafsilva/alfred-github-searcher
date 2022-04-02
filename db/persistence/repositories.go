package persistence

import (
	aw "github.com/deanishe/awgo"
)

type Repository struct {
	Url            string `json:"url"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	OwnerImagePath string `json:"ownerImagePath"`
}

func NewRepositoryJsonFile(wf *aw.Workflow) *JsonFile[Repository] {
	return &JsonFile[Repository]{
		dataType: "repositories",
		wf:       wf,
	}
}
