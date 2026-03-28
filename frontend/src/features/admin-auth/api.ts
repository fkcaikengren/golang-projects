import { adminHttp } from '../../shared/lib/http'
import type { AdminAuthPayload, AdminDashboardPayload, ApiEnvelope } from '../../shared/types/api'

export interface AdminLoginInput {
  email: string
  password: string
}

export async function adminLogin(input: AdminLoginInput): Promise<AdminAuthPayload> {
  const { data } = await adminHttp.post<ApiEnvelope<AdminAuthPayload>>('/login', input)
  return data.data
}

export async function fetchAdminDashboard(): Promise<AdminDashboardPayload> {
  const { data } = await adminHttp.get<ApiEnvelope<AdminDashboardPayload>>('')
  return data.data
}
