package cmd

import (
	"flag"
	"fmt"
	"os"

	"gist-downloader/internal/git"
	"gist-downloader/internal/github"
	"gist-downloader/internal/parser"
	"gist-downloader/internal/progress"
)

type Config struct {
	Token   string
	Output  string
	DryRun  bool
	Workers int
	URL     string
}

func Execute() {
	config := parseFlags()

	if config.URL == "" {
		fmt.Fprintln(os.Stderr, "Error: user URL is required")
		fmt.Fprintln(os.Stderr, "Usage: gist-downloader <user-url> [options]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	parsedURL, err := parser.ParseUserURL(config.URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing URL: %v\n", err)
		os.Exit(1)
	}

	client := github.NewClient(config.Token)

	fmt.Printf("Fetching gists for user: %s\n", parsedURL.Username)

	gists, err := client.ListGists(parsedURL.Username)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching gists: %v\n", err)
		os.Exit(1)
	}

	if len(gists) == 0 {
		fmt.Println("No gists found for this user")
		os.Exit(0)
	}

	fmt.Printf("Found %d gists\n", len(gists))

	tracker := progress.NewTracker()
	tracker.SetDryRun(config.DryRun)
	cloner := git.NewCloner(config.DryRun)
	executor := git.NewExecutor(cloner, tracker, config.Workers, config.Output, parsedURL.Username)

	if err := executor.Execute(gists); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.Token, "token", "", "GitHub personal access token (optional)")
	flag.StringVar(&config.Output, "output", "./gist", "Output directory for gists")
	flag.BoolVar(&config.DryRun, "dry-run", false, "Print operations without executing")
	flag.IntVar(&config.Workers, "workers", 5, "Number of parallel download workers")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: gist-downloader <user-url> [options]\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	config.URL = flag.Arg(0)

	return config
}
