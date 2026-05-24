package github

import "github.com/diagnosis/deploywatchv2/internal/config"

type GitHubClient struct {
	cfg *config.Config
}

func NewGitHubClient(cfg *config.Config) *GitHubClient {
	return &GitHubClient{cfg}
}
