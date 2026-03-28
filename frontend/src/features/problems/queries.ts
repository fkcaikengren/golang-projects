import { useMutation, useQuery } from '@tanstack/vue-query'

import { getProblem, listProblems, submitProblem } from './api'

export function useProblemsQuery() {
  return useQuery({
    queryKey: ['problems'],
    queryFn: listProblems,
  })
}

export function useProblemDetailQuery(slug: string) {
  return useQuery({
    queryKey: ['problems', slug],
    queryFn: () => getProblem(slug),
  })
}

export function useSubmitProblemMutation() {
  return useMutation({
    mutationFn: submitProblem,
  })
}
