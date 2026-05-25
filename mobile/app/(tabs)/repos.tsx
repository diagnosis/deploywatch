import { View, Text, FlatList, Pressable, ActivityIndicator, Alert, Modal, ScrollView } from 'react-native'
import { SafeAreaView } from 'react-native-safe-area-context'
import { Ionicons } from '@expo/vector-icons'
import { colors } from '@/src/styles/colors'
import { useWatchedRepos, useUnwatchRepo, useWatchRepo, useGitHubRepos } from '@/src/hooks/useRepos'
import { useState } from 'react'

export default function ReposScreen() {
    const { data: repos, isPending, refetch, isRefetching } = useWatchedRepos()
    const unwatchRepo = useUnwatchRepo()
    const watchRepo = useWatchRepo()
    const [showModal, setShowModal] = useState(false)
    const { data: githubRepos, isPending: loadingGithub } = useGitHubRepos()

    function handleUnwatch(id: number, name: string) {
        Alert.alert('Unwatch', `Stop watching ${name}?`, [
            { text: 'Cancel', style: 'cancel' },
            { text: 'Unwatch', style: 'destructive', onPress: () => unwatchRepo.mutate(id) }
        ])
    }

    return (
        <SafeAreaView style={{ flex: 1, backgroundColor: colors.background }}>
            {/* Header */}
            <View style={{ flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between', paddingHorizontal: 16, paddingVertical: 12, borderBottomWidth: 1, borderBottomColor: colors.border }}>
                <Text style={{ color: colors.textPrimary, fontSize: 20, fontWeight: '700' }}>Watched Repos</Text>
                <Pressable
                    onPress={() => setShowModal(true)}
                    style={{ width: 32, height: 32, borderRadius: 8, backgroundColor: colors.primaryDim, borderWidth: 1, borderColor: colors.borderAlt, justifyContent: 'center', alignItems: 'center' }}
                >
                    <Ionicons name="add" size={18} color={colors.primary} />
                </Pressable>
            </View>

            {isPending ? (
                <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center' }}>
                    <ActivityIndicator color={colors.primary} />
                </View>
            ) : (
                <FlatList
                    data={repos ?? []}
                    keyExtractor={item => String(item.repo_id)}
                    renderItem={({ item }) => (
                        <View style={{ flexDirection: 'row', alignItems: 'center', paddingHorizontal: 16, paddingVertical: 14, borderBottomWidth: 1, borderBottomColor: colors.border }}>
                            <Ionicons name="git-branch-outline" size={16} color={colors.primary} style={{ marginRight: 12 }} />
                            <Text style={{ color: colors.textPrimary, fontSize: 14, fontFamily: 'monospace', flex: 1 }}>{item.repo_full_name}</Text>
                            <Pressable onPress={() => handleUnwatch(item.repo_id, item.repo_full_name)}>
                                <Ionicons name="trash-outline" size={16} color={colors.error} />
                            </Pressable>
                        </View>
                    )}
                    ListEmptyComponent={
                        <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center', paddingTop: 80 }}>
                            <Ionicons name="git-branch-outline" size={48} color={colors.textDim} />
                            <Text style={{ color: colors.textSecondary, marginTop: 16 }}>No repos watched</Text>
                            <Pressable onPress={() => setShowModal(true)} style={{ marginTop: 12 }}>
                                <Text style={{ color: colors.primary }}>+ Watch a repo</Text>
                            </Pressable>
                        </View>
                    }
                    refreshing={isRefetching}
                    onRefresh={refetch}
                />
            )}

            {/* Add repo modal */}
            <Modal visible={showModal} animationType="slide" transparent onRequestClose={() => setShowModal(false)}>
                <View style={{ flex: 1, backgroundColor: 'rgba(0,0,0,0.6)', justifyContent: 'flex-end' }}>
                    <View style={{ backgroundColor: colors.surface, borderTopLeftRadius: 20, borderTopRightRadius: 20, padding: 20, maxHeight: '70%' }}>
                        <View style={{ flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
                            <Text style={{ color: colors.textPrimary, fontSize: 18, fontWeight: '700' }}>Watch a Repo</Text>
                            <Pressable onPress={() => setShowModal(false)}>
                                <Ionicons name="close" size={24} color={colors.textSecondary} />
                            </Pressable>
                        </View>
                        {loadingGithub ? (
                            <ActivityIndicator color={colors.primary} />
                        ) : (
                            <ScrollView>
                                {(githubRepos ?? []).map(repo => {
                                    const alreadyWatched = repos?.some(r => r.repo_id === repo.id)
                                    return (
                                        <Pressable
                                            key={repo.id}
                                            onPress={() => {
                                                if (alreadyWatched) return
                                                watchRepo.mutate({ repoId: repo.id, repoFullName: repo.full_name, installationId: repo.installation_id })
                                                setShowModal(false)
                                            }}
                                            style={{ flexDirection: 'row', alignItems: 'center', paddingVertical: 12, borderBottomWidth: 1, borderBottomColor: colors.border, opacity: alreadyWatched ? 0.4 : 1 }}
                                        >
                                            <Ionicons name="git-branch-outline" size={16} color={colors.primary} style={{ marginRight: 12 }} />
                                            <Text style={{ color: colors.textPrimary, fontSize: 14, fontFamily: 'monospace', flex: 1 }}>{repo.full_name}</Text>
                                            {alreadyWatched && <Text style={{ color: colors.textDim, fontSize: 12 }}>watching</Text>}
                                        </Pressable>
                                    )
                                })}
                            </ScrollView>
                        )}
                    </View>
                </View>
            </Modal>
        </SafeAreaView>
    )
}