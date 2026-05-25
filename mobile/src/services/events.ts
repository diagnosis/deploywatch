import { apiClient } from './api'
import { EventsResponse } from '@/src/types'

export async function getEvents(repoId?: number, page = 1, limit = 25): Promise<EventsResponse> {
    const params = new URLSearchParams({ page: String(page), limit: String(limit) })
    if (repoId) params.set('repo_id', String(repoId))
    return apiClient.get<EventsResponse>(`/api/events?${params.toString()}`)
}