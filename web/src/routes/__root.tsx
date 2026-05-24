// src/routes/__root.tsx

import {createRootRouteWithContext, Outlet} from "@tanstack/react-router";
import type {QueryClient} from "@tanstack/react-query";
import {NotFound} from "@/components/app/NotFound.tsx";


export const Route = createRootRouteWithContext<{
    queryClient: QueryClient
}>()({
    component: RootLayout,
    notFoundComponent: NotFound
})

function RootLayout() {
    return (
        <div>
            <Outlet />
        </div>
    )
}