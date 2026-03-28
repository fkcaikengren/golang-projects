import { useQuery } from '@tanstack/vue-query'

import { listSubmissions } from './api'

export function useSubmissionsQuery() {
  return useQuery({
    queryKey: ['submissions'],
    queryFn: listSubmissions,
  })
}
