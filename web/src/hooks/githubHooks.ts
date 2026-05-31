import { useQuery } from "@tanstack/react-query";
import { githubService } from "@/services/githubService.ts";

export const githubKeys = {
    repos: ['github', 'repos'] as const,
}

export const useGithubRepos = () => useQuery({
    queryKey: githubKeys.repos,
    queryFn: async () => {
        const res = await githubService.listRepos()
        if (!res.ok) return []
        return res.data
    },
    staleTime: 5 * 60 * 1000,
})