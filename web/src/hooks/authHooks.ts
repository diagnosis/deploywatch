import {authService} from "@/services/authService.ts";
import {useMutation, useQueryClient} from "@tanstack/react-query";
import {useRouter} from "@tanstack/react-router";


export const authKeys = {
    me: ['auth', 'me'] as const
}

export const meQueryOptions = () =>({
    queryKey: authKeys.me,
    queryFn: async () => {
        const res = await authService.me()
        if (!res.ok) throw new Error('Not authenticated')
        return res.data
    },
    staleTime: 5 * 60 * 1000,
    retry: false,
})

export const useLogout = () => {
    const qc = useQueryClient()
    const router = useRouter()
    return useMutation({
        mutationFn: () => authService.logout(),
        onSuccess: () => {
            qc.clear()
            router.navigate({ to: '/login' })
        },
    })
}