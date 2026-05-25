import * as SecureStore from 'expo-secure-store'
import {env} from "../config/env";


const ACCESS_TOKEN_KEY = 'access_token'
const REFRESH_TOKEN_KEY = 'refresh_token'

export const tokenStore = {
    getAccessToken: () => SecureStore.getItemAsync(ACCESS_TOKEN_KEY),
    getRefreshToken: () => SecureStore.getItemAsync(REFRESH_TOKEN_KEY),
    setTokens: async (access: string, refresh: string) => {
        await Promise.all([
            SecureStore.setItemAsync(ACCESS_TOKEN_KEY, access),
            SecureStore.setItemAsync(REFRESH_TOKEN_KEY, refresh),
        ])
    },
    clear: async () => {
        await Promise.all([
            SecureStore.deleteItemAsync(ACCESS_TOKEN_KEY),
            SecureStore.deleteItemAsync(REFRESH_TOKEN_KEY),
        ])
    },
}

let refreshPromise: Promise<string> | null = null

async function refreshAccessToken(): Promise<string> {
    const refreshToken = await tokenStore.getRefreshToken()
    if (!refreshToken) throw new Error('No refresh token')

    const res = await fetch(`${env.API_BASE_URL}/api/auth/refresh/mobile`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: refreshToken }),
    })

    if (!res.ok) {
        await tokenStore.clear()
        throw new Error('Refresh failed')
    }

    const json = await res.json()
    const { access_token, refresh_token } = json.data
    await tokenStore.setTokens(access_token, refresh_token)
    return access_token
}

async function request<T>(endpoint: string, init: RequestInit = {}): Promise<T> {
    const url = `${env.API_BASE_URL}${endpoint}`

    const doFetch = async () => {
        const token = await tokenStore.getAccessToken()
        const headers = new Headers(init.headers)
        headers.set('Content-Type', 'application/json')
        if (token) headers.set('Authorization', `Bearer ${token}`)
        return fetch(url, { ...init, headers })
    }

    let res = await doFetch()

    if (res.status === 401) {
        if (!refreshPromise) {
            refreshPromise = refreshAccessToken().finally(() => { refreshPromise = null })
        }
        await refreshPromise
        res = await doFetch()
    }

    if (!res.ok) {
        const json = await res.json().catch(() => null)
        throw new Error(json?.error?.message || res.statusText || 'Request failed')
    }

    const json = await res.json().catch(() => null)
    return (json?.data ?? json) as T
}

export const apiClient = {
    get: <T>(endpoint: string) => request<T>(endpoint),
    post: <T>(endpoint: string, body?: unknown) =>
        request<T>(endpoint, { method: 'POST', body: body ? JSON.stringify(body) : undefined }),
    delete: <T>(endpoint: string) => request<T>(endpoint, { method: 'DELETE' }),
}