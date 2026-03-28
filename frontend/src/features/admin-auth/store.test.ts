import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'

import { useAdminAuthStore } from './store'

describe('admin auth store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('starts anonymous, sets session, and clears session', () => {
    const store = useAdminAuthStore()

    expect(store.isAuthenticated).toBe(false)
    expect(store.token).toBe('')
    expect(store.adminUser).toBeNull()

    store.setSession({
      token: 'admin-token',
      admin_user: {
        id: 1,
        email: 'admin@go-oj.dev',
        display_name: 'Super Admin',
        status: 'active',
        last_login_at: 0,
      },
    })

    expect(store.isAuthenticated).toBe(true)
    expect(store.token).toBe('admin-token')
    expect(store.adminUser?.email).toBe('admin@go-oj.dev')

    store.clearSession()

    expect(store.isAuthenticated).toBe(false)
    expect(store.token).toBe('')
    expect(store.adminUser).toBeNull()
  })
})
