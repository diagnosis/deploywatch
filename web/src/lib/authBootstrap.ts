// src/lib/authBootstrap.ts
import { apiClient } from '@/lib/apiClient'
import type { User } from '@/types/auth'

export async function authBootstrap(): Promise<User | null> {
    const res = await apiClient.get<User>('/api/auth/me')
    if (res.ok) return res.data
    return null
}