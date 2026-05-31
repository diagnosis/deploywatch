import { useQuery } from "@tanstack/react-query";
import { githubService } from "@/services/githubService.ts";

export const githubKeys = {
    repos: ['github', 'repos'] as const,
    installation: ['github', 'installation'] as const,
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

export const useHasInstallation = () => useQuery({
    queryKey: githubKeys.installation,
    queryFn: async () => {
        const res = await githubService.checkInstallation()
        return res.ok ? res.data.installed : false
    },
    refetchOnWindowFocus: true,
})