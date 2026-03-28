import { env } from '../../shared/config/env'
import { http } from '../../shared/lib/http'
import { mockListSubmissions } from '../../shared/mocks/server'
import type { SubmissionItem } from '../../shared/types/api'

export async function listSubmissions(): Promise<SubmissionItem[]> {
  if (env.useMock) {
    return mockListSubmissions()
  }

  const { data } = await http.get<SubmissionItem[]>('/submissions')
  return data
}
