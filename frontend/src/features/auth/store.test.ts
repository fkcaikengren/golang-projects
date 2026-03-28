import { setActivePinia, createPinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'

import { useAuthStore } from './store'

describe('auth store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('starts anonymous, sets session, and clears session', () => {
    const store = useAuthStore()

    expect(store.isAuthenticated).toBe(false)
    expect(store.token).toBe('')
    expect(store.user).toBeNull()

    store.setSession({
      token: 'token-1',
      user: {
        id: 1,
        email: 'tester@example.com',
        nickname: 'tester',
        status: 'active',
      },
    })

    expect(store.isAuthenticated).toBe(true)
    expect(store.token).toBe('token-1')
    expect(store.user?.email).toBe('tester@example.com')

    store.clearSession()

    expect(store.isAuthenticated).toBe(false)
    expect(store.token).toBe('')
    expect(store.user).toBeNull()
  })
})
