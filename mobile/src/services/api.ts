import * as SecureStore from 'expo-secure-store'
import {env} from "../config/env";

const ACCESS_TOKEN_KEY = "access_token"
const REFRESH_TOKEN_KEY = "refresh_token"

export const tokenStore = {
    getAccessToken: () => SecureStore.getItemAsync(ACCESS_TOKEN_KEY),
    getRefreshToken: () => SecureStore.getItemAsync(REFRESH_TOKEN_KEY),
    setTokens: async (access:string, refresh:string) => {
        await Promise.all([
            SecureStore.setItemAsync(ACCESS_TOKEN_KEY, access),
            SecureStore.setItemAsync(REFRESH_TOKEN_KEY, refresh)
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

async function refreshAccessToken():Promise<string>{
    const refreshToken = await tokenStore.getRefreshToken()
    if (!refreshToken) throw new Error("No refresh token available")

    const res = await  fetch(`${env.API_BASE_URL}/api/auth/refresh`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({refresh_token:refreshToken})
    })

    if (!res.ok){
        await tokenStore.clear()
        throw new Error("Refresh token request failed")
    }

    const json = await res.json()
    const {access_token, refresh_token} = json.data
    await tokenStore.setTokens(access_token, refresh_token)
    return access_token
}