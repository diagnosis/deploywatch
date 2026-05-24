//src/lib/api/apiClient.ts

const BASE_URL = import.meta.env.VITE_API_URL ?? `http://localhost:8080`

type ApiError = {message: string, status?: number}
type ApiResponse<T> = {ok: true; data: T} | {ok: false; error: ApiError; status: number}

let refreshPromise: Promise<boolean> | null = null

async function refreshAccessToken(){
    if (refreshPromise) return refreshPromise
    refreshPromise = (async ()=>{
        try {
            const res = await fetch(`${BASE_URL}/api/auth/refresh`, {
                credentials: 'include',
                method: 'POST',
                headers: {'Content-Type': 'application/json'}
            })
            return res.ok;
        }finally {
            refreshPromise = null
        }
    })()

    return refreshPromise
}

async function request<T>(
    endpoint: string,
    init: RequestInit & { skipAuthRefresh?: boolean } = {},
) : Promise<ApiResponse<T>>{
    const url = `${BASE_URL}${endpoint}`

    const doFetch = () => {
        const headers = new Headers(init.headers)
        return fetch(url, {...init, headers, credentials: 'include'})
    }

    const res = await doFetch()
    const status = res.status

    const json = await res.json().catch(() => null)
    const payload = json?.data ?? json
    if (res.ok) return {ok: true, data: payload as T}

    const isRefreshCall = endpoint.startsWith('/api/auth/refresh')
    const canRetry = status === 401 && !init.skipAuthRefresh && !isRefreshCall

    if (canRetry){
        const refreshed = await refreshAccessToken()
        if (refreshed) {
            const res2 = await doFetch()
            const status2 =  res2.status
            const json2 = await res2.json().catch(() => null)
            const payload2 = json2?.data ?? json2
            if (res2.ok) return {ok: true, data: payload2 as T}

            return {
                ok: false,
                status: status2,
                error: { status: status2, message: json2?.error?.message || res2.statusText || 'Request failed' },
            }

        }
    }
    return {
        ok: false,
        status,
        error: { status, message: json?.error?.message || res.statusText || 'Request failed' },
    }

}

export const apiClient = {
    get : <T>(endpoint:string, init?: RequestInit & { skipAuthRefresh?: boolean}) =>
        (request<T>(endpoint, {
            ...init,
            method: 'GET',
        })),
    post : <T>(endpoint:string, body?: unknown, init?: RequestInit & { skipAuthRefresh?: boolean}) =>
        (request<T>(endpoint, {
            ...init,
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: body == null ? undefined : JSON.stringify(body),
        })),
    del: <T>(endpoint:string, init?: RequestInit & { skipAuthRefresh?: boolean}) =>
        (request<T>(endpoint, {
            ...init,
            method: 'DELETE',
        })),
    patch: <T>(endpoint:string, body?: unknown, init?: RequestInit & { skipAuthRefresh?: boolean}) =>
        (request<T>(endpoint, {
            ...init,
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
            },
            body: body == null ? undefined : JSON.stringify(body),
        })),

}