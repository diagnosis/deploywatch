import {apiClient} from "@/lib/apiClient.ts";


export const authService = {
    me: () => apiClient.get('/api/auth/me'),
    logout: () => apiClient.post('/api/auth/logout'),
}