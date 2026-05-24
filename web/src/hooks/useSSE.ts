// src/hooks/useSSE.ts

import { useEffect } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { eventKeys } from '@/hooks/eventHooks'
import { BASE_URL } from '@/lib/config'



export function useSSE(){
    const qc = useQueryClient()

    useEffect(() => {
        const es = new EventSource(`${BASE_URL}/api/sse`, {withCredentials: true})

        es.onmessage = (e) => {
            console.log('SSE event received:', e.data)
            qc.invalidateQueries({queryKey: eventKeys.all})
        }
        const eventTypes = ['push', 'pull_request', 'pull_request_review', 'create', 'delete', 'installation', 'check_run', 'workflow_job']

        eventTypes.forEach(type => {
            es.addEventListener(type, () => {
                console.log('SSE event received:', type)
                qc.invalidateQueries({ queryKey: ['events'] })
            })
        })

        es.onerror = (err) => {
            console.error("SSE error:", err)
        }

        return () => {
            es.close()
        }


    }, []);
}