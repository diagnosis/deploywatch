package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Repo struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	FullName       string `json:"full_name"`
	InstallationID int64  `json:"installation_id"`
}

type reposResponse struct {
	Repositories []Repo `json:"repositories"`
}

func (g *GitHubClient) ListInstallationRepos(ctx context.Context, token string, installationID int64) ([]Repo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/installation/repositories?per_page=200", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "deploywatch")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list repos: status %d", resp.StatusCode)
	}

	var rr reposResponse
	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		return nil, err
	}
	for i := range rr.Repositories {
		rr.Repositories[i].InstallationID = installationID
	}
	return rr.Repositories, nil
}
