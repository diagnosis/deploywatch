import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getWatchedRepos, watchRepo, unwatchRepo, getGitHubRepos } from '@/src/services/repos'

export function useWatchedRepos() {
    return useQuery({
        queryKey: ['repos'],
        queryFn: getWatchedRepos,
        staleTime: 1000 * 60,
    })
}

export function useGitHubRepos() {
    return useQuery({
        queryKey: ['github-repos'],
        queryFn: getGitHubRepos,
        staleTime: 1000 * 60 * 5,
    })
}

export function useWatchRepo() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: ({ repoId, repoFullName, installationId }: { repoId: number; repoFullName: string; installationId: number }) =>
            watchRepo(repoId, repoFullName, installationId),
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['repos'] }),
    })
}

export function useUnwatchRepo() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (id: number) => unwatchRepo(id),
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['repos'] }),
    })
}