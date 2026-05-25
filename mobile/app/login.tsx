import { View, Text, Pressable, ActivityIndicator } from 'react-native'
import { useRouter } from 'expo-router'
import { useState } from 'react'
import { loginWithGithub } from '@/src/services/auth'
import { colors } from '@/src/styles/colors'
import { Ionicons } from '@expo/vector-icons'

export default function LoginScreen() {
    const router = useRouter()
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)

    async function handleLogin() {
        setError(null)
        setLoading(true)
        try {
            const success = await loginWithGithub()
            if (success) {
                router.replace('/(tabs)/feed')
            } else {
                setError('Login cancelled or failed. Try again.')
            }
        } catch (e: any) {
            setError(e.message || 'Something went wrong')
        } finally {
            setLoading(false)
        }
    }

    return (
        <View style={{ flex: 1, backgroundColor: colors.background, justifyContent: 'center', alignItems: 'center', padding: 32 }}>
            {/* Logo */}
            <View style={{ width: 64, height: 64, borderRadius: 16, backgroundColor: colors.primaryDim, borderWidth: 1, borderColor: colors.borderAlt, justifyContent: 'center', alignItems: 'center', marginBottom: 24 }}>
                <Ionicons name="flash" size={28} color={colors.primary} />
            </View>

            <Text style={{ color: colors.textPrimary, fontSize: 28, fontWeight: '700', marginBottom: 8 }}>deploywatch</Text>
            <Text style={{ color: colors.textSecondary, fontSize: 15, marginBottom: 48, textAlign: 'center' }}>
                Monitor your GitHub repos in real time
            </Text>

            {error && (
                <Text style={{ color: colors.error, fontSize: 14, marginBottom: 16, textAlign: 'center' }}>{error}</Text>
            )}

            <Pressable
                onPress={handleLogin}
                disabled={loading}
                style={({ pressed }) => ({
                    flexDirection: 'row',
                    alignItems: 'center',
                    gap: 10,
                    backgroundColor: colors.primary,
                    paddingVertical: 16,
                    paddingHorizontal: 32,
                    borderRadius: 12,
                    opacity: pressed || loading ? 0.7 : 1,
                    width: '100%',
                    justifyContent: 'center',
                })}
            >
                {loading ? (
                    <ActivityIndicator color={colors.background} />
                ) : (
                    <>
                        <Ionicons name="logo-github" size={20} color={colors.background} />
                        <Text style={{ color: colors.background, fontSize: 16, fontWeight: '700' }}>
                            Continue with GitHub
                        </Text>
                    </>
                )}
            </Pressable>
        </View>
    )
}