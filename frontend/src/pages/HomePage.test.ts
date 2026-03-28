import { mount } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query'
import { describe, expect, it } from 'vitest'

import HomePage from './HomePage.vue'

describe('HomePage', () => {
  it('renders the homepage hero and featured problem sets', async () => {
    const wrapper = mount(HomePage, {
      global: {
        plugins: [
          createPinia(),
          [VueQueryPlugin, { queryClient: new QueryClient() }],
        ],
      },
    })

    await new Promise((resolve) => setTimeout(resolve, 0))

    expect(wrapper.text()).toContain('开始刷题')
    expect(wrapper.text()).toContain('热门题单')
  })
})
