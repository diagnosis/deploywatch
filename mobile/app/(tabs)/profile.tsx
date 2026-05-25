import { View, Text, Pressable, Alert, ScrollView } from 'react-native'
import { SafeAreaView } from 'react-native-safe-area-context'
import { useRouter } from 'expo-router'
import { Ionicons } from '@expo/vector-icons'
import { colors } from '@/src/styles/colors'
import { logout } from '@/src/services/auth'
import { useQuery } from '@tanstack/react-query'
import { apiClient } from '@/src/services/api'
import type { User } from '@/src/types'

function useMe() {
    return useQuery({
        queryKey: ['me'],
        queryFn: async () => {
            const result = await apiClient.get<any>('/api/auth/me')
            return result.user as User
        },
        staleTime: 1000 * 60 * 5,
    })
}

export default function ProfileScreen() {
    const router = useRouter()
    const { data: user } = useMe()

    function handleLogout() {
        Alert.alert('Sign out', 'Are you sure?', [
            { text: 'Cancel', style: 'cancel' },
            {
                text: 'Sign out', style: 'destructive', onPress: async () => {
                    await logout()
                    router.replace('/login')
                }
            }
        ])
    }

    return (
        <SafeAreaView style={{ flex: 1, backgroundColor: colors.background }}>
            <ScrollView>
                {/* Header */}
                <View style={{ paddingHorizontal: 16, paddingVertical: 12, borderBottomWidth: 1, borderBottomColor: colors.border }}>
                    <Text style={{ color: colors.textPrimary, fontSize: 20, fontWeight: '700' }}>Profile</Text>
                </View>

                {/* Avatar */}
                <View style={{ alignItems: 'center', paddingVertical: 32 }}>
                    <View style={{ width: 80, height: 80, borderRadius: 40, backgroundColor: colors.primaryDim, borderWidth: 2, borderColor: colors.borderAlt, justifyContent: 'center', alignItems: 'center', marginBottom: 12 }}>
                        <Ionicons name="person" size={36} color={colors.primary} />
                    </View>
                    <Text style={{ color: colors.textPrimary, fontSize: 18, fontWeight: '700' }}>{user?.login ?? '—'}</Text>
                    {user?.name && <Text style={{ color: colors.textSecondary, fontSize: 14, marginTop: 4 }}>{user.name}</Text>}
                    {user?.email && <Text style={{ color: colors.textDim, fontSize: 13, marginTop: 2 }}>{user.email}</Text>}
                </View>

                {/* Info */}
                <View style={{ marginHorizontal: 16, backgroundColor: colors.surface, borderRadius: 12, borderWidth: 1, borderColor: colors.border, marginBottom: 24 }}>
                    <View style={{ flexDirection: 'row', alignItems: 'center', gap: 12, padding: 16, borderBottomWidth: 1, borderBottomColor: colors.border }}>
                        <Ionicons name="logo-github" size={18} color={colors.textSecondary} />
                        <Text style={{ color: colors.textSecondary, fontSize: 14, flex: 1 }}>GitHub</Text>
                        <Text style={{ color: colors.textPrimary, fontSize: 14 }}>@{user?.login}</Text>
                    </View>
                    {user?.email && (
                        <View style={{ flexDirection: 'row', alignItems: 'center', gap: 12, padding: 16 }}>
                            <Ionicons name="mail-outline" size={18} color={colors.textSecondary} />
                            <Text style={{ color: colors.textSecondary, fontSize: 14, flex: 1 }}>Email</Text>
                            <Text style={{ color: colors.textPrimary, fontSize: 14 }}>{user.email}</Text>
                        </View>
                    )}
                </View>

                {/* Sign out */}
                <Pressable
                    onPress={handleLogout}
                    style={({ pressed }) => ({ flexDirection: 'row', alignItems: 'center', justifyContent: 'center', gap: 8, marginHorizontal: 16, padding: 16, borderRadius: 12, backgroundColor: colors.error, opacity: pressed ? 0.8 : 1 })}
                >
                    <Ionicons name="log-out-outline" size={18} color="#fff" />
                    <Text style={{ color: '#fff', fontSize: 16, fontWeight: '600' }}>Sign out</Text>
                </Pressable>

                <Text style={{ color: colors.textDim, fontSize: 12, textAlign: 'center', marginTop: 24 }}>deploywatch v1.0.0</Text>
            </ScrollView>
        </SafeAreaView>
    )
}