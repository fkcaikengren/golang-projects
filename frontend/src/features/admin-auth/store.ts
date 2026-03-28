import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

import type { AdminAuthPayload, AdminUserProfile } from '../../shared/types/api'

export const useAdminAuthStore = defineStore('admin-auth', () => {
  const token = ref('')
  const adminUser = ref<AdminUserProfile | null>(null)

  const isAuthenticated = computed(() => token.value.length > 0 && adminUser.value !== null)

  function setSession(payload: AdminAuthPayload) {
    token.value = payload.token
    adminUser.value = payload.admin_user
  }

  function clearSession() {
    token.value = ''
    adminUser.value = null
  }

  return {
    token,
    adminUser,
    isAuthenticated,
    setSession,
    clearSession,
  }
})
