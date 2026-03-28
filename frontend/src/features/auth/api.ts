import { env } from '../../shared/config/env'
import { http } from '../../shared/lib/http'
import { mockLogin, mockRegister } from '../../shared/mocks/server'
import type { AuthPayload } from '../../shared/types/api'

export interface AuthInput {
  email: string
  password: string
  nickname?: string
}

export async function login(input: AuthInput): Promise<AuthPayload> {
  if (env.useMock) {
    return mockLogin()
  }

  const { data } = await http.post<AuthPayload>('/auth/login', input)
  return data
}

export async function register(input: AuthInput): Promise<AuthPayload> {
  if (env.useMock) {
    return mockRegister()
  }

  const { data } = await http.post<AuthPayload>('/auth/register', input)
  return data
}
