import type { App } from 'vue'

import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query'
import { createPinia } from 'pinia'

import { router } from '../router'

const queryClient = new QueryClient()

export function installAppProviders(app: App<Element>) {
  app.use(createPinia())
  app.use(router)
  app.use(VueQueryPlugin, { queryClient })
}

export { queryClient }
