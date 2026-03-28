import { describe, expect, it } from 'vitest'

import { createAppRouter } from './index'

describe('router', () => {
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
      ]),
    )
  })
})
