package github

type GistFile struct {
	Filename string `json:"filename"`
}

type GistResponse struct {
	ID          string              `json:"id"`
	Description string              `json:"description"`
	Public      bool                `json:"public"`
	GitPullURL  string              `json:"git_pull_url"`
	Files       map[string]GistFile `json:"files"`
}

type Gist struct {
	ID         string
	Name       string
	GitPullURL string
	Public     bool
}
