import {apiClient} from "@/lib/apiClient.ts";
import type {watchedRepo, watchedRepoReq, WatchedReposResponse} from "@/types/repo.ts";



export const repoService = {
    add : (req: watchedRepoReq) => (
        apiClient.post<watchedRepo>('/api/repos/watch', req)
    ),
    list: () => (
        apiClient.get<WatchedReposResponse>('/api/repos')
    ),
    remove: (repoID: number) =>
        apiClient.del(`/api/repos/watch/${repoID}`)
}