import axios from 'axios'

import { env } from '../config/env'
import { useAuthStore } from '../../features/auth/store'
import { useAdminAuthStore } from '../../features/admin-auth/store'

export const http = axios.create({
  baseURL: env.apiBaseUrl,
  timeout: 10000,
})

// 请求拦截器：自动添加 token
http.interceptors.request.use((config) => {
  const authStore = useAuthStore()
  if (authStore.token) {
    config.headers.Authorization = `Bearer ${authStore.token}`
  }
  return config
})

export const adminHttp = axios.create({
  baseURL: env.adminApiBaseUrl,
  timeout: 10000,
})

adminHttp.interceptors.request.use((config) => {
  const adminAuthStore = useAdminAuthStore()
  if (adminAuthStore.token) {
    config.headers.Authorization = `Bearer ${adminAuthStore.token}`
  }
  return config
})
