package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Cloner struct {
	dryRun bool
}

func NewCloner(dryRun bool) *Cloner {
	return &Cloner{
		dryRun: dryRun,
	}
}

func (c *Cloner) CloneGist(gitURL, gistID, name, basePath string) error {
	sanitized := sanitizeFilename(name)
	gistPath := filepath.Join(basePath, sanitized)

	if _, err := os.Stat(gistPath); err == nil {
		return fmt.Errorf("already exists")
	}

	if c.dryRun {
		fmt.Printf("[DRY-RUN] Would clone: %s -> %s\n", gitURL, gistPath)
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(gistPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	cmd := exec.Command("git", "clone", gitURL, gistPath)
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	return nil
}

func sanitizeFilename(name string) string {
	sanitized := strings.Map(func(r rune) rune {
		if r == '/' || r == '\\' || r == ':' || r == '*' || r == '?' || r == '"' || r == '<' || r == '>' || r == '|' {
			return '-'
		}
		return r
	}, name)

	sanitized = strings.TrimSpace(sanitized)
	if sanitized == "" {
		return "unnamed"
	}

	if len(sanitized) > 255 {
		sanitized = sanitized[:255]
	}

	return sanitized
}
