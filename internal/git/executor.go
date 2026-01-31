package git

import (
	"fmt"
	"sync"

	"gist-downloader/internal/github"
	"gist-downloader/internal/progress"
)

type Executor struct {
	cloner   *Cloner
	tracker  *progress.Tracker
	workers  int
	basePath string
	user     string
}

func NewExecutor(cloner *Cloner, tracker *progress.Tracker, workers int, basePath, user string) *Executor {
	return &Executor{
		cloner:   cloner,
		tracker:  tracker,
		workers:  workers,
		basePath: basePath,
		user:     user,
	}
}

func (e *Executor) Execute(gists []github.Gist) error {
	e.tracker.SetTotal(len(gists))

	gistChan := make(chan github.Gist, len(gists))
	for _, gist := range gists {
		gistChan <- gist
	}
	close(gistChan)

	var wg sync.WaitGroup
	wg.Add(e.workers)

	for i := 0; i < e.workers; i++ {
		go func() {
			defer wg.Done()
			e.worker(gistChan)
		}()
	}

	wg.Wait()
	e.tracker.Finalize()

	return nil
}

func (e *Executor) worker(gistChan <-chan github.Gist) {
	for gist := range gistChan {
		e.tracker.IncrementPending()

		userPath := fmt.Sprintf("%s/%s", e.basePath, e.user)
		err := e.cloner.CloneGist(gist.GitPullURL, gist.ID, gist.Name, userPath)

		if err != nil {
			if err.Error() == "already exists" {
				e.tracker.IncrementSkipped()
			} else {
				e.tracker.IncrementFailed(gist.ID, err)
			}
		} else {
			e.tracker.IncrementCompleted()
		}
	}
}
