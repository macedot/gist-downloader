package progress

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Tracker struct {
	total     int
	completed int
	failed    int
	skipped   int
	mu        sync.Mutex
	startTime time.Time
	failures  map[string]string
	dryRun    bool
}

func NewTracker() *Tracker {
	return &Tracker{
		failures: make(map[string]string),
	}
}

func (t *Tracker) SetDryRun(dryRun bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.dryRun = dryRun
}

func (t *Tracker) SetTotal(total int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.total = total
	t.startTime = time.Now()
	t.printHeader()
}

func (t *Tracker) IncrementPending() {
	t.printProgress()
}

func (t *Tracker) IncrementCompleted() {
	t.mu.Lock()
	t.completed++
	t.mu.Unlock()
	t.printProgress()
}

func (t *Tracker) IncrementFailed(id string, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.failed++
	t.failures[id] = err.Error()
}

func (t *Tracker) IncrementSkipped() {
	t.mu.Lock()
	t.skipped++
	t.mu.Unlock()
	t.printProgress()
}

func (t *Tracker) Finalize() {
	t.mu.Lock()
	total := t.total
	completed := t.completed
	failed := t.failed
	skipped := t.skipped
	duration := time.Since(t.startTime)
	failures := make(map[string]string)
	for k, v := range t.failures {
		failures[k] = v
	}
	t.mu.Unlock()

	fmt.Println()
	fmt.Printf("Completed in %v\n", duration.Round(time.Millisecond))
	fmt.Printf("Total: %d | Success: %d | Failed: %d | Skipped: %d\n", total, completed, failed, skipped)

	if len(failures) > 0 {
		fmt.Println("\nFailed gists:")
		for id, err := range failures {
			fmt.Printf("  %s: %s\n", id, err)
		}
	}
}

func (t *Tracker) printHeader() {
	fmt.Printf("Downloading %d gists...\n", t.total)
	if t.dryRun {
		fmt.Println("⚠️  DRY RUN MODE: No files will be downloaded")
		fmt.Println()
	}
}

func (t *Tracker) printProgress() {
	t.mu.Lock()
	defer t.mu.Unlock()

	done := t.completed + t.failed + t.skipped
	if t.total == 0 {
		return
	}

	percent := int(float64(done) / float64(t.total) * 100)
	barWidth := 40
	filled := int(float64(barWidth) * float64(done) / float64(t.total))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	fmt.Printf("\r[%s] %d%% | Done: %d/%d | Failed: %d | Skipped: %d", bar, percent, done, t.total, t.failed, t.skipped)
}
