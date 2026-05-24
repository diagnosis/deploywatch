import {apiClient} from "@/lib/apiClient.ts";
import type { EventsResponse } from "@/types/event.ts";


export const eventService = {
    list: (repoID?: number) => {
        const params = repoID ? `?repo_id=${repoID}` : ''
        return apiClient.get<EventsResponse>(`/api/events${params}`)
    }
}