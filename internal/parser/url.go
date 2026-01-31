package parser

import (
	"fmt"
	"net/url"
	"strings"
)

type UserURL struct {
	Username string
}

func ParseUserURL(rawURL string) (*UserURL, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Host != "gist.github.com" {
		return nil, fmt.Errorf("expected gist.github.com host, got %s", parsedURL.Host)
	}

	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(pathParts) < 1 || pathParts[0] == "" {
		return nil, fmt.Errorf("could not extract username from URL")
	}

	username := pathParts[0]

	return &UserURL{
		Username: username,
	}, nil
}
