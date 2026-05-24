//src/services/githubService.ts
import {apiClient} from "@/lib/apiClient.ts";

export interface GithubRepo{
    id: number;
    name: string;
    full_name: string;
    installation_id: number;
}

export const githubService = {
    listRepos: () => apiClient.get<GithubRepo[]>('/api/github/repos')
}