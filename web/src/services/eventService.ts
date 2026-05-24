import {apiClient} from "@/lib/apiClient.ts";
import type { EventsResponse } from "@/types/event.ts";


export const eventService = {
    list: (repoID?: number, page = 1, limit = 25) => {
        const params = new URLSearchParams()
        if (repoID) params.set("repo_id", String(repoID))
        params.set('page', String(page))
        params.set('limit', String(limit))
        return apiClient.get<EventsResponse>(`/api/events?${params.toString()}`)
    }
}