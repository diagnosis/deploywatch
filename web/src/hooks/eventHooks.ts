// src/hooks/eventHooks.ts

import {eventService} from "@/services/eventService.ts";

export const eventKeys = {
    all: ['events']
}

export const eventQueryOptions = (repoID?: number, page = 1) => ({
    queryKey: ['events', repoID, page],
    queryFn: async () => {
        const res = await eventService.list(repoID, page)
        if (!res.ok) throw new Error('Failed to fetch events')
        return res.data
    },
})