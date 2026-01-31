# gist-downloader

A Go CLI tool to download all gists from a GitHub user using `git clone`.

## Features

- Downloads all gists (public and private with token)
- Parallel downloads with configurable worker count
- Progress tracking with real-time updates
- Dry-run mode to preview operations
- Skip existing gists
- Organizes gists as `gist/<user>/<gist-name>/`

## Installation

```bash
go build -o gist-downloader .
```

## Usage

```bash
# Basic usage
./gist-downloader https://gist.github.com/username

# With authentication (higher rate limits and private gists)
./gist-downloader https://gist.github.com/username --token $GITHUB_TOKEN

# Custom output directory
./gist-downloader https://gist.github.com/username --output ./my-gists

# Dry run (preview what would be downloaded)
./gist-downloader https://gist.github.com/username --dry-run

# More workers for faster downloads
./gist-downloader https://gist.github.com/username --workers 10
```

## Options

- `--token`: GitHub personal access token (optional)
- `--output`: Output directory (default: `./gist`)
- `--dry-run`: Print operations without executing
- `--workers`: Number of parallel workers (default: 5)

## Requirements

- Go 1.16 or later
- Git installed and available in PATH

## How it works

1. Parses the GitHub user URL to extract the username
2. Fetches all gists via the GitHub API
3. Clones each gist repository using `git clone`
4. Organizes them in `gist/<user>/<gist-name>/` structure
5. Shows progress and final summary

## Rate Limits

- Without token: 60 requests/hour
- With token: 5000 requests/hour

For frequent use, create a GitHub token at https://github.com/settings/tokens
