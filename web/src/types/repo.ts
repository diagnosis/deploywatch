// src/repo.ts

export interface watchedRepo {
    id: string,
    user_id: string,
    repo_id: number,
    repo_full_name: string,
    event_filters : Record<string, unknown>,
    created_at: Date,
    updated_at?: Date,
}
export interface watchedRepoReq {
    repo_id: number,
    repo_full_name: string,
    installation_id:number
}

export interface WatchedReposResponse {
    repos: watchedRepo[] | null
    count: number
}