package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	token     string
	client    *http.Client
	userAgent string
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent: "gist-downloader",
	}
}

func (c *Client) ListGists(username string) ([]Gist, error) {
	var allGists []Gist
	page := 1

	for {
		reqURL := fmt.Sprintf("https://api.github.com/users/%s/gists?page=%d&per_page=100", username, page)
		req, err := http.NewRequest("GET", reqURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("User-Agent", c.userAgent)
		if c.token != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
		}

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch gists: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				return nil, fmt.Errorf("user '%s' not found", username)
			}
			if resp.StatusCode == http.StatusUnauthorized {
				return nil, fmt.Errorf("unauthorized: invalid token")
			}
			if resp.StatusCode == http.StatusForbidden {
				return nil, fmt.Errorf("rate limit exceeded. Please provide a GitHub token with --token")
			}
			return nil, fmt.Errorf("API error: %s", resp.Status)
		}

		var gists []GistResponse
		if err := json.NewDecoder(resp.Body).Decode(&gists); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		resp.Body.Close()

		if len(gists) == 0 {
			break
		}

		for _, gist := range gists {
			name := gist.Description
			if name == "" {
				for _, file := range gist.Files {
					name = file.Filename
					break
				}
			}
			if name == "" {
				name = gist.ID
			}

			allGists = append(allGists, Gist{
				ID:         gist.ID,
				Name:       name,
				GitPullURL: gist.GitPullURL,
				Public:     gist.Public,
			})
		}

		page++
	}

	return allGists, nil
}
