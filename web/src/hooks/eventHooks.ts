// src/hooks/eventHooks.ts

import {eventService} from "@/services/eventService.ts";

export const eventKeys = {
    all: ['events']
}

export const eventQueryOptions = (repoID: number | undefined) => ({
    queryKey: ['events', repoID],
    queryFn: async () => {
      const res =  await eventService.list(repoID)
        if (res.ok){
            return res.data.events ?? []
        }
        throw new Error("list empty")
    },

});