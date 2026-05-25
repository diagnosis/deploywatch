import * as WebBrowser from 'expo-web-browser'
import * as Linking from 'expo-linking'
import { env } from '@/src/config/env'
import { tokenStore } from '@/src/services/api'

export async function loginWithGithub(): Promise<boolean> {
    const redirectUri = Linking.createURL('/auth/callback')
    const result = await WebBrowser.openAuthSessionAsync(
        `${env.API_BASE_URL}/api/auth/github/login?mobile=true&redirect_uri=${encodeURIComponent(redirectUri)}`,
        redirectUri
    )

    if (result.type !== 'success') {
        return false
    }

    const url = result.url
    const params = new URL(url).searchParams
    const accessToken = params.get('access_token')
    const refreshToken = params.get('refresh_token')

    if (!accessToken || !refreshToken) return false

    await tokenStore.setTokens(accessToken, refreshToken)
    return true
}

export async function logout(): Promise<void> {
    await tokenStore.clear()
}