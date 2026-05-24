// src/types/auth.ts
export interface User {
    id: string,
    github_id: number,
    login: string,
    name?: string,
    avatar_url?: string,
    email?: string,
    created_at?: Date,
    updated_at?: Date,
}