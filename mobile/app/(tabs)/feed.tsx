import { View, Text, FlatList, RefreshControl, Pressable, ActivityIndicator } from 'react-native'
import { SafeAreaView } from 'react-native-safe-area-context'
import { useState } from 'react'
import { useEvents } from '@/src/hooks/useEvents'
import { useWatchedRepos } from '@/src/hooks/useRepos'
import { Ionicons } from '@expo/vector-icons'
import { colors } from '@/src/styles/colors'
import type { Event } from '@/src/types'
import * as Linking from 'expo-linking'

function timeAgo(date: string) {
    const seconds = Math.floor((Date.now() - new Date(date).getTime()) / 1000)
    if (seconds < 60) return `${seconds}s ago`
    if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`
    if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`
    return `${Math.floor(seconds / 86400)}d ago`
}

function eventColor(type: string) {
    switch (type) {
        case 'push': return colors.push
        case 'pull_request': return colors.pullRequest
        case 'pull_request_review': return colors.review
        default: return colors.textSecondary
    }
}

function getEventURL(event: Event, repoFullName: string): string {
    let p: any = {}
    try {
        p = JSON.parse(atob(event.payload))
    } catch {}
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
        <Pressable
            onPress={() => Linking.openURL(getEventURL(event, repoFullName))}
            style={({ pressed }) => ({
                flexDirection: 'row',
                alignItems: 'flex-start',
                gap: 12,
                paddingHorizontal: 16,
                paddingVertical: 14,
                borderBottomWidth: 1,
                borderBottomColor: colors.border,
                backgroundColor: pressed ? colors.primaryDim : 'transparent',
            })}
        >
            <View style={{ width: 32, height: 32, borderRadius: 8, backgroundColor: colors.surface, borderWidth: 1, borderColor: colors.border, justifyContent: 'center', alignItems: 'center' }}>
                <Ionicons
                    name={event.event_type === 'push' ? 'git-commit' : event.event_type === 'pull_request' ? 'git-pull-request' : 'flash'}
                    size={14}
                    color={eventColor(event.event_type)}
                />
            </View>
            <View style={{ flex: 1 }}>
                <View style={{ flexDirection: 'row', alignItems: 'center', gap: 8, flexWrap: 'wrap' }}>
                    <Text style={{ color: colors.textPrimary, fontSize: 14, fontWeight: '600' }}>{event.event_type}</Text>
                    {event.action && (
                        <View style={{ backgroundColor: colors.primaryDim, borderRadius: 10, paddingHorizontal: 8, paddingVertical: 2, borderWidth: 1, borderColor: colors.borderAlt }}>
                            <Text style={{ color: colors.primary, fontSize: 11 }}>{event.action}</Text>
                        </View>
                    )}
                </View>
                <View style={{ flexDirection: 'row', alignItems: 'center', gap: 6, marginTop: 2 }}>
                    <Ionicons name="person-outline" size={11} color={colors.textDim} />
                    <Text style={{ color: colors.textSecondary, fontSize: 12 }}>{event.actor_login}</Text>
                    <Text style={{ color: colors.textDim, fontSize: 12 }}>·</Text>
                    <Text style={{ color: colors.textDim, fontSize: 12 }}>{timeAgo(event.received_at)}</Text>
                </View>
            </View>
        </Pressable>
    )
}

export default function FeedScreen() {
    const [selectedRepo, setSelectedRepo] = useState<number | undefined>(undefined)
    const [page, setPage] = useState(1)
    const { data: repos } = useWatchedRepos()
    const { data, isPending, isRefetching, refetch } = useEvents(selectedRepo, page)
    const repoName = repos?.find(r => r.repo_id === selectedRepo)?.repo_full_name

    return (
        <SafeAreaView style={{ flex: 1, backgroundColor: colors.background }}>
            {/* Header */}
            <View style={{ flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between', paddingHorizontal: 16, paddingVertical: 12, borderBottomWidth: 1, borderBottomColor: colors.border }}>
                <View style={{ flexDirection: 'row', alignItems: 'center', gap: 8 }}>
                    <Ionicons name="flash" size={18} color={colors.primary} />
                    <Text style={{ color: colors.textPrimary, fontSize: 16, fontWeight: '700' }}>deploywatch</Text>
                </View>
                <View style={{ flexDirection: 'row', alignItems: 'center', gap: 6, backgroundColor: colors.primaryDim, borderRadius: 10, paddingHorizontal: 8, paddingVertical: 4, borderWidth: 1, borderColor: colors.borderAlt }}>
                    <View style={{ width: 6, height: 6, borderRadius: 3, backgroundColor: colors.primary }} />
                    <Text style={{ color: colors.primary, fontSize: 11 }}>live</Text>
                </View>
            </View>

            {/* Repo filter */}
            <FlatList
                horizontal
                data={[{ id: -1, repo_id: -1, repo_full_name: 'All repos', created_at: '' }, ...(repos ?? [])]}
                keyExtractor={item => String(item.repo_id)}
                showsHorizontalScrollIndicator={false}
                contentContainerStyle={{ paddingHorizontal: 12, paddingVertical: 8, gap: 8 , alignItems:'center'}}
                renderItem={({ item }) => {
                    const isSelected = item.repo_id === -1 ? selectedRepo === undefined : selectedRepo === item.repo_id
                    return (
                        <Pressable
                            onPress={() => {
                                setSelectedRepo(item.repo_id === -1 ? undefined : item.repo_id)
                                setPage(1)
                            }}
                            style={{
                                paddingHorizontal: 12, paddingVertical: 6, borderRadius: 8,
                                backgroundColor: isSelected ? colors.primaryDim : colors.surface,
                                borderWidth: 1, borderColor: isSelected ? colors.primary : colors.border,
                            }}
                        >
                            <Text style={{ color: isSelected ? colors.primary : colors.textSecondary, fontSize: 12, fontFamily: 'monospace' }}>
                                {item.repo_full_name}
                            </Text>
                        </Pressable>
                    )
                }}
            />

            {/* Events */}
            {isPending ? (
                <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center' }}>
                    <ActivityIndicator color={colors.primary} />
                </View>
            ) : (
                <FlatList
                    data={data?.events ?? []}
                    keyExtractor={item => String(item.id)}
                    renderItem={({ item }) => (
                        <EventCard
                            event={item}
                            repoFullName={repos?.find(r => r.repo_id === item.repo_id)?.repo_full_name ?? ''}
                        />
                    )}
                    refreshControl={
                        <RefreshControl refreshing={isRefetching} onRefresh={refetch} tintColor={colors.primary} />
                    }
                    ListEmptyComponent={
                        <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center', paddingTop: 80 }}>
                            <Ionicons name="flash-outline" size={48} color={colors.textDim} />
                            <Text style={{ color: colors.textSecondary, marginTop: 16 }}>No events yet</Text>
                        </View>
                    }
                    ListFooterComponent={
                        data && data.total_pages > 1 ? (
                            <View style={{ flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', padding: 16 }}>
                                <Pressable onPress={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1}>
                                    <Text style={{ color: page === 1 ? colors.textDim : colors.primary }}>← Prev</Text>
                                </Pressable>
                                <Text style={{ color: colors.textSecondary, fontSize: 12 }}>{page} / {data.total_pages}</Text>
                                <Pressable onPress={() => setPage(p => Math.min(data.total_pages, p + 1))} disabled={page === data.total_pages}>
                                    <Text style={{ color: page === data.total_pages ? colors.textDim : colors.primary }}>Next →</Text>
                                </Pressable>
                            </View>
                        ) : null
                    }
                />
            )}
        </SafeAreaView>
    )
}