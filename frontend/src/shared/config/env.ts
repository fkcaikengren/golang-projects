const rawApiBaseUrl = import.meta.env.VITE_API_BASE_URL
const rawUseMock = import.meta.env.VITE_USE_MOCK

export const env = {
  apiBaseUrl:
    typeof rawApiBaseUrl === 'string' && rawApiBaseUrl.length > 0
      ? rawApiBaseUrl
      : '/api/v1',
  useMock: rawUseMock !== 'false',
}
