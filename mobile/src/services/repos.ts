import { apiClient } from './api'
import {GitHubRepo, WatchedRepo} from '@/src/types'

export async function getWatchedRepos(): Promise<WatchedRepo[]> {
    const result = await apiClient.get<any>('/api/repos')
    return result?.repos ?? []
}

export async function watchRepo(repoId: number, repoFullName: string, installationId: number): Promise<void> {
    return apiClient.post('/api/repos/watch', {
        repo_id: repoId,
        repo_full_name: repoFullName,
        installation_id: installationId
    })
}

export async function unwatchRepo(id: number): Promise<void> {
    return apiClient.delete(`/api/repos/watch/${id}`)
}

export async function getGitHubRepos(): Promise<GitHubRepo[]> {
    const result = await apiClient.get<any>('/api/github/repos')
    return result
}