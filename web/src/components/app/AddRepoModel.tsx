// src/components/app/AddRepoModal.tsx
import { useAddRepo } from "@/hooks/repoHooks.ts"
import {useEffect, useState} from "react"
import { FaXmark, FaCodeBranch } from 'react-icons/fa6'
import {type GithubRepo, githubService} from "@/services/githubService.ts";

export function AddRepoModal({ onClose }: { onClose: () => void }) {
    const addRepo = useAddRepo()
    const [repos, setRepos] = useState<GithubRepo[]>([])
    const [selectedRepo, setSelectedRepo] = useState<GithubRepo | null>(null)
    const [search, setSearch] = useState<string>("")

    useEffect(() => {
        githubService.listRepos().then(res => {
            if (res.ok){
                setRepos(res.data)
            }else{
                console.error("failed to fetch repos", res.error.message)
            }

        })
    }, [])

    const handleSubmit = async (e: React.SubmitEvent) => {
        e.preventDefault()
        if (!selectedRepo) return
        const res = await addRepo.mutateAsync({
            repo_id: selectedRepo.id,
            repo_full_name: selectedRepo.full_name,
            installation_id: selectedRepo.installation_id,
        })
        if (res.ok) onClose()
    }

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
            <div className="absolute inset-0 bg-black/70 backdrop-blur-sm" onClick={onClose} />

            <div className="relative z-10 w-full max-w-sm mx-4 bg-[#0a1520] border border-cyan-500/[0.12] rounded-2xl p-6 shadow-2xl shadow-cyan-500/5">

                {/* Glow */}
                <div className="absolute -top-px left-1/2 -translate-x-1/2 w-32 h-px bg-gradient-to-r from-transparent via-cyan-500/40 to-transparent" />

                <div className="flex items-center justify-between mb-5">
                    <div className="flex items-center gap-2">
                        <div className="w-6 h-6 rounded-md bg-cyan-500/10 border border-cyan-500/20 flex items-center justify-center">
                            <FaCodeBranch size={11} className="text-cyan-400/70" />
                        </div>
                        <h2 className="text-cyan-100/70 text-sm font-medium">Watch a repo</h2>
                    </div>
                    <button onClick={onClose} className="text-white/20 hover:text-cyan-400/50 transition-colors">
                        <FaXmark size={14} />
                    </button>
                </div>

                <form onSubmit={handleSubmit} className="flex flex-col gap-3">
                    {/* Search input */}
                    <input
                        placeholder="Search repos..."
                        value={search}
                        onChange={e => setSearch(e.target.value)}
                        className="w-full bg-cyan-500/[0.04] border border-cyan-500/[0.10] rounded-lg px-3 py-2 text-cyan-100/70 text-sm placeholder-white/15 outline-none focus:border-cyan-500/40 transition-all font-mono"
                    />

                    {/* Repo list */}
                    <div className="max-h-48 overflow-y-auto flex flex-col gap-1">
                        {repos
                            .filter(r => r.full_name.toLowerCase().includes(search.toLowerCase()))
                            .map(repo => (
                                <div
                                    key={repo.id}
                                    onClick={() => setSelectedRepo(repo)}
                                    className={`px-3 py-2 rounded-lg cursor-pointer text-sm font-mono transition-colors ${
                                        selectedRepo?.id === repo.id
                                            ? 'bg-cyan-500/20 text-cyan-300 border border-cyan-500/30'
                                            : 'text-white/40 hover:bg-cyan-500/[0.06] hover:text-white/70'
                                    }`}
                                >
                                    {repo.full_name}
                                </div>
                            ))}
                    </div>

                    <div className="flex gap-2 mt-2">
                        <button
                            type="button"
                            onClick={onClose}
                            className="flex-1 py-2 rounded-lg border border-white/[0.06] text-white/30 text-sm hover:text-white/50 hover:bg-white/[0.02] transition-colors"
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            disabled={addRepo.isPending}
                            className="flex-1 py-2 rounded-lg bg-cyan-500/20 hover:bg-cyan-500/30 border border-cyan-500/30 hover:border-cyan-500/50 text-cyan-300 text-sm font-medium transition-all disabled:opacity-40"
                        >
                            {addRepo.isPending ? 'Adding...' : 'Watch repo'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    )
}