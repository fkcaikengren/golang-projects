import { config } from '@vue/test-utils'

config.global.stubs = {
  transition: false,
  teleport: true,
  RouterLink: {
    template: '<a><slot /></a>',
  },
}

Object.defineProperty(window, 'scrollTo', {
  writable: true,
  value: () => {},
})
