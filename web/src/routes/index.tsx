// src/routes/index.tsx
import { createFileRoute } from '@tanstack/react-router'
import { requireAuth } from "@/lib/routeGuard.ts"
import { useSSE } from "@/hooks/useSSE.ts"
import { useSuspenseQuery } from "@tanstack/react-query"
import { reposQueryOptions, useRemoveRepo } from "@/hooks/repoHooks.ts"
import { useEffect, useRef, useState } from "react"
import { eventQueryOptions } from "@/hooks/eventHooks.ts"
import { AddRepoModal } from "@/components/app/AddRepoModel.tsx"
import { useLogout } from "@/hooks/authHooks.ts"
import { FaGithub, FaCodeBranch, FaTrash, FaPlus, FaBolt, FaCodePullRequest, FaComment } from 'react-icons/fa6'
import type { Event } from "@/types/event.ts"
import {useGithubRepos, useHasInstallation} from "@/hooks/githubHooks.ts"

const GITHUB_APP_INSTALL_URL = 'https://github.com/apps/deploywatch/installations/new'

export const Route = createFileRoute('/')({
    beforeLoad: async ({ context }) => requireAuth(context.queryClient),
    component: RouteComponent,
})

function eventIcon(type: string) {
    switch (type) {
        case 'push': return <FaCodeBranch className="text-cyan-400" size={13} />
        case 'pull_request': return <FaCodePullRequest className="text-teal-400" size={13} />
        case 'pull_request_review': return <FaComment className="text-sky-400" size={13} />
        default: return <FaBolt className="text-cyan-500/40" size={13} />
    }
}

function eventBadgeColor(type: string) {
    switch (type) {
        case 'push': return 'bg-cyan-500/10 text-cyan-400/70 border-cyan-500/20'
        case 'pull_request': return 'bg-teal-500/10 text-teal-400/70 border-teal-500/20'
        case 'pull_request_review': return 'bg-sky-500/10 text-sky-400/70 border-sky-500/20'
        default: return 'bg-white/[0.04] text-white/30 border-white/[0.06]'
    }
}

function timeAgo(date: string | Date) {
    const seconds = Math.floor((Date.now() - new Date(date).getTime()) / 1000)
    if (seconds < 60) return `${seconds}s ago`
    if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`
    if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`
    return `${Math.floor(seconds / 86400)}d ago`
}

function getEventURL(event: Event, repoFullName: string): string {
    const raw = event.payload
    let p: any = {}
    try {
        let jsonStr: string
        if (typeof raw === 'string') {
            try {
                jsonStr = atob(raw)
            } catch {
                jsonStr = raw
            }
            p = JSON.parse(jsonStr)
        }
    } catch {
        return `https://github.com/${repoFullName}`
    }
    switch (event.event_type) {
        case 'push': return `https://github.com/${repoFullName}/commit/${p.after}`
        case 'pull_request': return `https://github.com/${repoFullName}/pull/${p.number}`
        case 'pull_request_review': return `https://github.com/${repoFullName}/pull/${p.pull_request?.number}`
        case 'create': return `https://github.com/${repoFullName}/tree/${p.ref}`
        default: return `https://github.com/${repoFullName}`
    }
}

function EventCard({ event, repoFullName }: { event: Event; repoFullName: string }) {
    return (
        <a href={getEventURL(event, repoFullName)} target="_blank" rel="noopener noreferrer">
            <div className="flex items-start gap-3 px-4 py-3.5 border-b border-white/[0.04] hover:bg-cyan-500/[0.02] transition-colors group">
                <div className="mt-0.5 w-7 h-7 rounded-lg bg-[#0d1f2d] border border-cyan-500/10 flex items-center justify-center flex-shrink-0 group-hover:border-cyan-500/20 transition-colors">
                    {eventIcon(event.event_type)}
                </div>
                <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 flex-wrap">
                        <span className="text-white/75 text-sm font-medium">{event.event_type}</span>
                        {event.action && (
                            <span className={`text-[10px] px-1.5 py-0.5 rounded-full border ${eventBadgeColor(event.event_type)}`}>
                                {event.action}
                            </span>
                        )}
                    </div>
                    <div className="flex items-center gap-1.5 mt-0.5">
                        <FaGithub size={10} className="text-white/20" />
                        <span className="text-white/30 text-xs">{event.actor_login}</span>
                        <span className="text-white/10 text-xs">·</span>
                        <span className="text-white/20 text-xs">{timeAgo(event.received_at)}</span>
                    </div>
                </div>
                <div className="w-1.5 h-1.5 rounded-full bg-cyan-500/30 mt-1.5 flex-shrink-0" />
            </div>
        </a>
    )
}

function RouteComponent() {
    useSSE()
    const [page, setPage] = useState(1)
    const { data: watchedRepos } = useSuspenseQuery(reposQueryOptions())
    const [selectedRepo, setSelectedRepo] = useState<number | undefined>(undefined)
    const events = useSuspenseQuery(eventQueryOptions(selectedRepo, page))
    const [showAddModal, setShowAddModal] = useState(false)
    const removeRepo = useRemoveRepo()
    const logout = useLogout()
    const prevRepoRef = useRef(selectedRepo)
    const selectedRepoName = watchedRepos.find(r => r.repo_id === selectedRepo)?.repo_full_name
    const { data: hasInstallation = false } = useHasInstallation()

    useEffect(() => {
        if (prevRepoRef.current !== selectedRepo) {
            prevRepoRef.current = selectedRepo
            if (page !== 1) setPage(1)
        }
    }, [page, selectedRepo])

    return (
        <div className="h-screen bg-[#080d12] flex flex-col">
            <div className="fixed top-0 left-1/2 -translate-x-1/2 w-[800px] h-[300px] bg-cyan-500/[0.03] blur-[100px] pointer-events-none" />

            <header className="h-12 border-b border-cyan-500/[0.08] flex items-center justify-between px-4 flex-shrink-0 relative z-10 bg-[#080d12]/80 backdrop-blur-sm">
                <div className="flex items-center gap-2">
                    <div className="w-6 h-6 rounded-md bg-cyan-500/20 border border-cyan-500/30 flex items-center justify-center">
                        <FaBolt size={11} className="text-cyan-400" />
                    </div>
                    <span className="text-cyan-100/80 text-sm font-medium tracking-tight">deploywatch</span>
                    <div className="flex items-center gap-1 ml-2 px-1.5 py-0.5 rounded-full bg-cyan-500/10 border border-cyan-500/20">
                        <div className="w-1 h-1 rounded-full bg-cyan-400 animate-pulse" />
                        <span className="text-cyan-400/70 text-[10px]">live</span>
                    </div>
                </div>
                <button
                    onClick={() => logout.mutate()}
                    className="text-white/20 hover:text-cyan-400/60 text-xs transition-colors px-2 py-1 rounded-md hover:bg-cyan-500/[0.06]"
                >
                    Sign out
                </button>
            </header>

            <div className="flex flex-1 overflow-hidden relative z-10">
                <aside className="w-56 border-r border-cyan-500/[0.08] flex flex-col flex-shrink-0 bg-[#080d12]">
                    <div className="flex items-center justify-between px-3 py-3 border-b border-cyan-500/[0.06]">
                        <span className="text-cyan-500/40 text-[10px] uppercase tracking-widest font-semibold">Repos</span>
                        <button
                            onClick={() => setShowAddModal(true)}
                            className="w-5 h-5 rounded-md bg-cyan-500/10 hover:bg-cyan-500/20 border border-cyan-500/20 hover:border-cyan-500/40 flex items-center justify-center transition-all"
                        >
                            <FaPlus size={8} className="text-cyan-400/70" />
                        </button>
                    </div>
                    <div className="flex-1 overflow-y-auto py-1">
                        {watchedRepos.length === 0 ? (
                            <div className="px-3 py-8 text-center">
                                <div className="w-8 h-8 rounded-xl bg-cyan-500/[0.06] border border-cyan-500/10 flex items-center justify-center mx-auto mb-3">
                                    <FaCodeBranch size={12} className="text-cyan-500/30" />
                                </div>
                                <p className="text-white/15 text-xs mb-3">No repos watched</p>
                                <button
                                    onClick={() => setShowAddModal(true)}
                                    className="text-cyan-400/60 hover:text-cyan-400 text-xs transition-colors"
                                >
                                    + Watch a repo
                                </button>
                            </div>
                        ) : (
                            <>
                                <button
                                    onClick={() => setSelectedRepo(undefined)}
                                    className={`w-full flex items-center gap-2 px-3 py-2 text-left transition-all ${
                                        selectedRepo === undefined
                                            ? 'bg-cyan-500/[0.08] text-cyan-300/80 border-r-2 border-cyan-500/50'
                                            : 'text-white/35 hover:text-cyan-300/60 hover:bg-cyan-500/[0.04]'
                                    }`}
                                >
                                    <FaBolt size={10} />
                                    <span className="text-xs font-medium">All repos</span>
                                </button>
                                {watchedRepos.map(repo => (
                                    <div
                                        key={repo.id}
                                        className={`group flex items-center gap-2 px-3 py-2 cursor-pointer transition-all ${
                                            selectedRepo === repo.repo_id
                                                ? 'bg-cyan-500/[0.08] text-cyan-300/80 border-r-2 border-cyan-500/50'
                                                : 'text-white/35 hover:text-cyan-300/60 hover:bg-cyan-500/[0.04]'
                                        }`}
                                        onClick={() => setSelectedRepo(repo.repo_id)}
                                    >
                                        <FaCodeBranch size={10} className="flex-shrink-0" />
                                        <span className="text-xs flex-1 truncate font-mono">{repo.repo_full_name}</span>
                                        <button
                                            onClick={e => { e.stopPropagation(); removeRepo.mutate(repo.repo_id) }}
                                            className="opacity-0 group-hover:opacity-100 transition-opacity text-white/15 hover:text-red-400/70"
                                        >
                                            <FaTrash size={10} />
                                        </button>
                                    </div>
                                ))}
                            </>
                        )}
                    </div>
                </aside>

                <main className="flex-1 flex flex-col overflow-hidden">
                    <div className="px-4 py-2.5 border-b border-cyan-500/[0.06] flex items-center gap-2 bg-[#080d12]/50">
                        <FaCodeBranch size={11} className="text-cyan-500/30" />
                        <span className="text-white/40 text-sm font-mono">{selectedRepoName ?? 'all repos'}</span>
                        <span className="text-cyan-500/40 text-xs ml-auto font-mono">{events.data?.total} events</span>
                    </div>

                    <div className="flex-1 overflow-y-auto">
                        {!events.data?.events || events.data.events.length === 0 ? (
                            <div className="flex flex-col items-center justify-center h-full text-center px-8">
                                {!hasInstallation ? (
                                    <div className="flex flex-col items-center">
                                        <div className="w-12 h-12 rounded-2xl bg-cyan-500/10 border border-cyan-500/20 flex items-center justify-center mb-4">
                                            <FaGithub size={20} className="text-cyan-400/60" />
                                        </div>
                                        <p className="text-white/50 text-sm font-medium mb-2">Install GitHub App first</p>
                                        <p className="text-white/20 text-xs mb-6">Connect your GitHub account to start monitoring repos</p>
                                        <a
                                            href={GITHUB_APP_INSTALL_URL}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            className="px-4 py-2 rounded-lg bg-cyan-500/20 border border-cyan-500/30 text-cyan-300 text-sm hover:bg-cyan-500/30 transition-colors"
                                        >
                                            Install GitHub App →
                                        </a>
                                    </div>
                                ) : (
                                    <div className="flex flex-col items-center">
                                        <FaBolt size={32} className="text-cyan-500/20 mb-4" />
                                        <p className="text-white/20 text-sm">No events yet</p>
                                        <p className="text-white/10 text-xs mt-1">Watch a repo to start monitoring</p>
                                    </div>
                                )}
                            </div>
                        ) : (
                            events.data.events.map(event => (
                                <EventCard
                                    key={event.id}
                                    event={event}
                                    repoFullName={watchedRepos.find(r => r.repo_id === event.repo_id)?.repo_full_name ?? ''}
                                />
                            ))
                        )}
                    </div>

                    {events.data?.total_pages > 1 && (
                        <div className="flex items-center justify-between px-4 py-3 border-t border-cyan-500/[0.06]">
                            <button
                                onClick={() => setPage(p => Math.max(1, p - 1))}
                                disabled={page === 1}
                                className="text-xs text-cyan-400/60 hover:text-cyan-400 disabled:opacity-20 transition-colors"
                            >
                                ← Prev
                            </button>
                            <span className="text-white/20 text-xs font-mono">{page} / {events.data.total_pages}</span>
                            <button
                                onClick={() => setPage(p => Math.min(events.data.total_pages, p + 1))}
                                disabled={page === events.data?.total_pages}
                                className="text-xs text-cyan-400/60 hover:text-cyan-400 disabled:opacity-20 transition-colors"
                            >
                                Next →
                            </button>
                        </div>
                    )}
                </main>
            </div>

            {showAddModal && <AddRepoModal onClose={() => setShowAddModal(false)} />}
            <footer className="h-8 border-t border-cyan-500/[0.06] flex items-center justify-between px-4 bg-[#080d12]">
                <span className="text-white/30 text-xs font-mono">deploywatch v1.0</span>
                <span className="text-white/30 text-xs">built by <span className="text-cyan-400/50">diagnosis</span></span>
            </footer>
        </div>
    )
}