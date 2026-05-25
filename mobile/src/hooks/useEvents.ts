import { useQuery } from '@tanstack/react-query'
import { getEvents } from '@/src/services/events'

export function useEvents(repoId?: number, page = 1) {
    return useQuery({
        queryKey: ['events', repoId, page],
        queryFn: () => getEvents(repoId, page),
        staleTime: 1000 * 30,
    })
}