package db

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os/exec"

	aw "github.com/deanishe/awgo"
)

func downloadImage(wf *aw.Workflow, url string) string {
	id := imageId(url)
	path := fmt.Sprintf("%s/%s", wf.CacheDir(), id)

	// check if it exists or it is already running
	if wf.Cache.Exists(id) {
		return path
	}
	jobName := "image-download/%s"
	if wf.IsRunning(jobName) {
		return path
	}

	// schedule the execution in background
	log.Printf("downloading image from %s..\n", url)
	cmd := exec.Command("curl", "-o", path, url)
	err := wf.RunInBackground(jobName, cmd)
	if err != nil {
		log.Printf("error while download image: %s\n", err.Error())
	}

	// the job will not complete, but that's fine
	return path
}

func imageId(url string) string {
	encoder := sha1.New()
	encoder.Write([]byte(url))
	hv := encoder.Sum(nil)
	return hex.EncodeToString(hv)
}
