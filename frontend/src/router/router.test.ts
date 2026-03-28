import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'

import { createAppRouter } from './index'
import { useAdminAuthStore } from '../features/admin-auth/store'

describe('router', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('registers the MVP routes', () => {
    const router = createAppRouter()

    const paths = router.getRoutes().map((route) => route.path)

    expect(paths).toEqual(
      expect.arrayContaining([
        '/',
        '/login',
        '/register',
        '/problem-sets',
        '/problem-sets/:slug',
        '/problems',
        '/problems/:slug',
        '/submissions',
        '/progress',
        '/admin/login',
        '/admin',
      ]),
    )
  })

  it('redirects anonymous admin visitors to /admin/login', async () => {
    const router = createAppRouter()

    await router.push('/admin')

    expect(router.currentRoute.value.fullPath).toBe('/admin/login')
  })

  it('allows authenticated admin visitors to access /admin', async () => {
    const router = createAppRouter()
    const store = useAdminAuthStore()
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

    await router.push('/admin')

    expect(router.currentRoute.value.fullPath).toBe('/admin')
  })
})
