import type {
  AuthPayload,
  ProblemDetail,
  ProblemSetSummary,
  ProblemSummary,
  SubmissionItem,
} from '../types/api'
import {
  mockAuth,
  mockProblemDetails,
  mockProblems,
  mockProblemSets,
  mockSubmissions,
} from './data'

function sleep(ms = 10) {
  return new Promise((resolve) => window.setTimeout(resolve, ms))
}

export async function mockLogin(): Promise<AuthPayload> {
  await sleep()
  return mockAuth
}

export async function mockRegister(): Promise<AuthPayload> {
  await sleep()
  return mockAuth
}

export async function mockListProblemSets(): Promise<ProblemSetSummary[]> {
  await sleep()
  return mockProblemSets
}

export async function mockGetProblemSet(slug: string): Promise<ProblemSetSummary> {
  await sleep()
  return mockProblemSets.find((item) => item.slug === slug) ?? mockProblemSets[0]
}

export async function mockListProblems(): Promise<ProblemSummary[]> {
  await sleep()
  return mockProblems
}

export async function mockGetProblem(slug: string): Promise<ProblemDetail> {
  await sleep()
  return mockProblemDetails[slug] ?? mockProblemDetails['two-sum']
}

export async function mockListSubmissions(): Promise<SubmissionItem[]> {
  await sleep()
  return mockSubmissions
}

export async function mockSubmitProblem(): Promise<{ message: string }> {
  await sleep()
  return { message: '提交成功，已进入占位判题流程。' }
}
