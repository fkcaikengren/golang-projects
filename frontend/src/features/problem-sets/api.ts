import { env } from '../../shared/config/env'
import { http } from '../../shared/lib/http'
import { mockGetProblemSet, mockListProblemSets } from '../../shared/mocks/server'
import type { ProblemSetSummary } from '../../shared/types/api'

export async function listProblemSets(): Promise<ProblemSetSummary[]> {
  if (env.useMock) {
    return mockListProblemSets()
  }

  const { data } = await http.get<ProblemSetSummary[]>('/problem-sets')
  return data
}

export async function getProblemSet(slug: string): Promise<ProblemSetSummary> {
  if (env.useMock) {
    return mockGetProblemSet(slug)
  }

  const { data } = await http.get<ProblemSetSummary>(`/problem-sets/${slug}`)
  return data
}
