package db

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	aw "github.com/deanishe/awgo"
)

type workflowScheduler struct {
	wf *aw.Workflow
}

func (s *workflowScheduler) scheduleRepositoriesRefresh() error {
	return s.scheduleRefresh("repository")
}

func (s *workflowScheduler) schedulePullRequestsRefresh() error {
	return s.scheduleRefresh("pr")
}

func (s *workflowScheduler) scheduleRefresh(target string) error {
	job := fmt.Sprintf("sync-%s", target)
	if s.wf.IsRunning(job) {
		log.Printf("there's already a %s job running\n", job)
		return nil
	}

	log.Printf("scheduling %s to be run in the background..\n", job)
	cmd := exec.Command(os.Args[0], "sync", target)
	if err := s.wf.RunInBackground(job, cmd); err != nil {
		log.Printf("error scheduling refresh: %s\n", err)
	}

	log.Printf("successfully scheduled %s job\n", job)
	return nil
}
