
import './index.css'
import  ReactDOM  from 'react-dom/client'
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import {createRouter, RouterProvider} from "@tanstack/react-router";
import {routeTree} from "./routeTree.gen";
import {PageSpinner} from "@/components/ui/PageSpinner.tsx";
import {authBootstrap} from "@/lib/authBootstrap.ts";
import {StrictMode} from "react";

const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            retry: 1,
            staleTime: 60 * 1000,
            refetchOnWindowFocus: false,
        }
    }
})

const router = createRouter({
    routeTree,
    defaultPendingComponent: PageSpinner,
    defaultPendingMs:0,
    defaultPendingMinMs:500,
    context: { queryClient },
})

declare module '@tanstack/react-router' {
    interface Register {
        router: typeof router
    }
}

const rootElement = document.getElementById('root')!

async function start() {
    // ⬇️ important: wait for refresh attempt
    await authBootstrap()

    ReactDOM.createRoot(rootElement).render(
        <StrictMode>
            <QueryClientProvider client={queryClient}>
                <RouterProvider router={router} />
            </QueryClientProvider>
        </StrictMode>,
    )
}

start()