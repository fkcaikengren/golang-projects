import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

import type { AuthPayload, UserProfile } from '../../shared/types/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref('')
  const user = ref<UserProfile | null>(null)

  const isAuthenticated = computed(() => token.value.length > 0 && user.value !== null)

  function setSession(payload: AuthPayload) {
    token.value = payload.token
    user.value = payload.user
  }

  function clearSession() {
    token.value = ''
    user.value = null
  }

  return {
    token,
    user,
    isAuthenticated,
    setSession,
    clearSession,
  }
})
