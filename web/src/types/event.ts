// src/types/event.ts
export interface Event {
    id: string,
    repo_id: number,
    event_type: string,
    action?: string,
    actor_login: string,
    payload: Record<string, unknown>,
    received_at: Date,
}

export interface EventsResponse {
    events: Event[] | null
    count: number
}