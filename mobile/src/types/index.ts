export interface User {
    id: string
    login: string
    name?: string
    avatar_url?: string
    email?: string
    git_hub_id?: number
}

export interface GitHubRepo {
    id: number
    name: string
    full_name: string
    installation_id: number
}


export interface WatchedRepo {
    id: number
    repo_id: number
    repo_full_name: string
    created_at: string
}

export interface Event {
    id: number
    repo_id: number
    event_type: string
    action?: string
    actor_login: string
    payload: string
    received_at: string
}

export interface EventsResponse {
    events: Event[]
    total: number
    total_pages: number
    page: number
    limit: number
}