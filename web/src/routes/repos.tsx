// src/routes/repos.tsx

import { createFileRoute } from '@tanstack/react-router'
import {requireAuth} from "@/lib/routeGuard.ts";

export const Route = createFileRoute('/repos')({
  beforeLoad: ({context}) => requireAuth(context.queryClient),
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/repos"!</div>
}
