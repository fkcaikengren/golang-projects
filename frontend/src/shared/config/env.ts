const rawApiBaseUrl = import.meta.env.VITE_API_BASE_URL
const rawAdminApiBaseUrl = import.meta.env.VITE_ADMIN_API_BASE_URL
const rawUseMock = import.meta.env.VITE_USE_MOCK

export const env = {
  apiBaseUrl:
    typeof rawApiBaseUrl === 'string' && rawApiBaseUrl.length > 0
      ? rawApiBaseUrl
      : '/api/v1',
  adminApiBaseUrl:
    typeof rawAdminApiBaseUrl === 'string' && rawAdminApiBaseUrl.length > 0
      ? rawAdminApiBaseUrl
      : import.meta.env.DEV
        ? 'http://127.0.0.1:8080/admin'
        : '/admin',
  useMock: rawUseMock !== 'false',
}
