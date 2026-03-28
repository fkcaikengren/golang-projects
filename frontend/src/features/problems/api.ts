import { env } from '../../shared/config/env'
import { http } from '../../shared/lib/http'
import { mockGetProblem, mockListProblems, mockSubmitProblem } from '../../shared/mocks/server'
import type { ProblemDetail, ProblemSummary } from '../../shared/types/api'

export async function listProblems(): Promise<ProblemSummary[]> {
  if (env.useMock) {
    return mockListProblems()
  }

  const { data } = await http.get<ProblemSummary[]>('/problems')
  return data
}

export async function getProblem(slug: string): Promise<ProblemDetail> {
  if (env.useMock) {
    return mockGetProblem(slug)
  }

  const { data } = await http.get<ProblemDetail>(`/problems/${slug}`)
  return data
}

export async function submitProblem(payload: { slug: string; code: string; language: string }) {
  if (env.useMock) {
    return mockSubmitProblem()
  }

  const { data } = await http.post('/submissions', payload)
  return data
}
