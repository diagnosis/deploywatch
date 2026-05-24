// src/routes/index.tsx
import { createFileRoute } from '@tanstack/react-router'
import { requireAuth } from "@/lib/routeGuard.ts"
import { useSSE } from "@/hooks/useSSE.ts"
import { useSuspenseQuery } from "@tanstack/react-query"
import { reposQueryOptions, useRemoveRepo } from "@/hooks/repoHooks.ts"
import { useState } from "react"
import { eventQueryOptions } from "@/hooks/eventHooks.ts"
import { AddRepoModal } from "@/components/app/AddRepoModel.tsx"
import { useLogout } from "@/hooks/authHooks.ts"
import { FaGithub, FaCodeBranch, FaTrash, FaPlus, FaBolt, FaCodePullRequest, FaComment } from 'react-icons/fa6'
import type { Event } from "@/types/event.ts"

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

function EventCard({ event }: { event: Event }) {
    return (
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
            {/* Live dot */}
            <div className="w-1.5 h-1.5 rounded-full bg-cyan-500/30 mt-1.5 flex-shrink-0" />
        </div>
    )
}

function RouteComponent() {
    useSSE()

    const { data: watchedRepos } = useSuspenseQuery(reposQueryOptions())
    const [selectedRepo, setSelectedRepo] = useState<number | undefined>(undefined)
    const events = useSuspenseQuery(eventQueryOptions(selectedRepo))
    const [showAddModal, setShowAddModal] = useState(false)
    const removeRepo = useRemoveRepo()
    const logout = useLogout()

    const selectedRepoName = watchedRepos.find(r => r.repo_id === selectedRepo)?.repo_full_name
    console.log('events',events)
    return (
        <div className="min-h-screen bg-[#080d12] flex flex-col">

            {/* Ambient glow */}
            <div className="fixed top-0 left-1/2 -translate-x-1/2 w-[800px] h-[300px] bg-cyan-500/[0.03] blur-[100px] pointer-events-none" />

            {/* Top bar */}
            <header className="h-12 border-b border-cyan-500/[0.08] flex items-center justify-between px-4 flex-shrink-0 relative z-10 bg-[#080d12]/80 backdrop-blur-sm">
                <div className="flex items-center gap-2">
                    <div className="w-6 h-6 rounded-md bg-cyan-500/20 border border-cyan-500/30 flex items-center justify-center">
                        <FaBolt size={11} className="text-cyan-400" />
                    </div>
                    <span className="text-cyan-100/80 text-sm font-medium tracking-tight">deploywatch</span>
                    {/* Live indicator */}
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

                {/* Sidebar */}
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
                                            <FaTrash size={8} />
                                        </button>
                                    </div>
                                ))}
                            </>
                        )}
                    </div>
                </aside>

                {/* Main */}
                <main className="flex-1 flex flex-col overflow-hidden">
                    <div className="px-4 py-2.5 border-b border-cyan-500/[0.06] flex items-center gap-2 bg-[#080d12]/50">
                        <FaCodeBranch size={11} className="text-cyan-500/30" />
                        <span className="text-white/40 text-sm font-mono">
                            {selectedRepoName ?? 'all repos'}
                        </span>
                        <span className="text-cyan-500/20 text-xs ml-auto font-mono">
                            {events.data?.length ?? 0} events
                        </span>
                    </div>

                    <div className="flex-1 overflow-y-auto">
                        {!events.data || events.data.length === 0 ? (
                            <div className="flex flex-col items-center justify-center h-full text-center">
                                <div className="w-12 h-12 rounded-2xl bg-cyan-500/[0.04] border border-cyan-500/10 flex items-center justify-center mb-4">
                                    <FaBolt size={18} className="text-cyan-500/20" />
                                </div>
                                <p className="text-white/20 text-sm">No events yet</p>
                                <p className="text-white/10 text-xs mt-1.5">Push to a watched repo to see live events</p>
                            </div>
                        ) : (
                            events.data.map(event => (
                                <EventCard key={event.id} event={event} />
                            ))
                        )}
                    </div>
                </main>
            </div>

            {showAddModal && <AddRepoModal onClose={() => setShowAddModal(false)} />}
        </div>
    )
}