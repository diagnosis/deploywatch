package handlers

import (
	"github.com/diagnosis/deploywatchv2/internal/github"
	"github.com/diagnosis/deploywatchv2/internal/store/installation"
)

type GitHubHandler struct {
	gitHubClient      *github.GitHubClient
	installationStore installation.InstallationStore
}

func NewGitHubHandler(client *github.GitHubClient, installationStore installation.InstallationStore) *GitHubHandler {
	return &GitHubHandler{
		gitHubClient:      client,
		installationStore: installationStore,
	}
}
