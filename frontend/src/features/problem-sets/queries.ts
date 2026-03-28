import { useQuery } from '@tanstack/vue-query'

import { getProblemSet, listProblemSets } from './api'

export function useProblemSetsQuery() {
  return useQuery({
    queryKey: ['problemSets'],
    queryFn: listProblemSets,
  })
}

export function useProblemSetDetailQuery(slug: string) {
  return useQuery({
    queryKey: ['problemSets', slug],
    queryFn: () => getProblemSet(slug),
  })
}
