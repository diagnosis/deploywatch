// src/hooks/repoHooks.ts

import {repoService} from "@/services/repoService.ts";
import {useMutation, useQueryClient} from "@tanstack/react-query";
import type {watchedRepoReq} from "@/types/repo.ts";

export const repoKeys = {
    all: ['repos'] as const
}

export const reposQueryOptions = () => ({
    queryKey: repoKeys.all,
    queryFn: async () => {
        const res = await repoService.list()
        if (!res.ok) throw new Error('Failed to fetch repos')
        return res.data.repos ?? []
    },
})

export const useAddRepo = () => {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (payload: watchedRepoReq)=>
            repoService.add(payload),
        onSuccess: () => qc.invalidateQueries({queryKey: repoKeys.all}),
    })
}

export const useRemoveRepo = () => {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (repoID: number)=>
            repoService.remove(repoID),
        onSuccess: () => qc.invalidateQueries({queryKey: repoKeys.all}),
    })
}