import { apiClient } from './api'
import { WatchedRepo } from '@/src/types'

export async function getWatchedRepos(): Promise<WatchedRepo[]> {
    return apiClient.get<WatchedRepo[]>('/api/repos')
}

export async function watchRepo(repoId: number, repoFullName: string): Promise<void> {
    return apiClient.post('/api/repos/watch', { repo_id: repoId, repo_full_name: repoFullName })
}

export async function unwatchRepo(id: number): Promise<void> {
    return apiClient.delete(`/api/repos/watch/${id}`)
}

export async function getGitHubRepos(): Promise<{ id: number; full_name: string }[]> {
    return apiClient.get('/api/github/repos')
}