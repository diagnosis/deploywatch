import type {QueryClient} from "@tanstack/react-query";
import {redirect} from "@tanstack/react-router";
import {meQueryOptions} from "@/hooks/authHooks.ts";

export async function requireAuth(queryClient: QueryClient){
    try {
        await queryClient.ensureQueryData(meQueryOptions())
    }catch {
        throw redirect({to: '/login'})
    }
}