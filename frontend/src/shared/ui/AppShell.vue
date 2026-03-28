<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { NButton, NLayout, NLayoutContent, NLayoutHeader, NSpace, NTag } from 'naive-ui'

import { useAuthStore } from '../../features/auth/store'

const route = useRoute()
const authStore = useAuthStore()

const currentPath = computed(() => route.path)
const navItems = [
  { label: '首页', to: '/' },
  { label: '题单', to: '/problem-sets' },
  { label: '题库', to: '/problems' },
  { label: '提交记录', to: '/submissions' },
]
</script>

<template>
  <n-layout position="absolute">
    <n-layout-header bordered class="shell-header">
      <div class="page-section shell-header-inner">
        <RouterLink class="brand" to="/">Go OJ</RouterLink>
        <n-space :size="20" align="center">
          <RouterLink
            v-for="item in navItems"
            :key="item.to"
            :class="['nav-link', { active: currentPath === item.to }]"
            :to="item.to"
          >
            {{ item.label }}
          </RouterLink>
        </n-space>
        <n-space align="center">
          <n-tag v-if="authStore.isAuthenticated" type="success" round>
            {{ authStore.user?.nickname }}
          </n-tag>
          <RouterLink v-if="!authStore.isAuthenticated" to="/login">
            <n-button type="primary" secondary>登录 / 注册</n-button>
          </RouterLink>
          <n-button v-else quaternary @click="authStore.clearSession()">退出</n-button>
        </n-space>
      </div>
    </n-layout-header>
    <n-layout-content embedded>
      <RouterView />
    </n-layout-content>
  </n-layout>
</template>

<style scoped>
.shell-header {
  background: rgba(255, 255, 255, 0.82);
  backdrop-filter: blur(12px);
}

.shell-header-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 72px;
}

.brand {
  font-size: 22px;
  font-weight: 800;
  letter-spacing: 0.04em;
}

.nav-link {
  color: #516079;
  transition: color 0.2s ease;
}

.nav-link.active,
.nav-link:hover {
  color: #1f5eff;
}
</style>
