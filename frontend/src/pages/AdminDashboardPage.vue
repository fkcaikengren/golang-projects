<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { AxiosError } from 'axios'
import { NButton, NCard, NGrid, NGridItem, NResult, NSkeleton, NSpace, NTag } from 'naive-ui'

import { fetchAdminDashboard } from '../features/admin-auth/api'
import { useAdminAuthStore } from '../features/admin-auth/store'
import type { AdminDashboardPayload } from '../shared/types/api'

const router = useRouter()
const adminAuthStore = useAdminAuthStore()

const pending = ref(true)
const errorMessage = ref('')
const dashboard = ref<AdminDashboardPayload | null>(null)

const welcomeTitle = computed(() => dashboard.value?.title ?? 'Admin Dashboard')

async function loadDashboard() {
  pending.value = true
  errorMessage.value = ''

  try {
    dashboard.value = await fetchAdminDashboard()
  } catch (error) {
    if (error instanceof AxiosError && error.response?.status === 401) {
      adminAuthStore.clearSession()
      await router.replace('/admin/login')
      return
    }
    errorMessage.value = '后台数据加载失败，请稍后刷新重试。'
  } finally {
    pending.value = false
  }
}

async function handleLogout() {
  adminAuthStore.clearSession()
  await router.push('/admin/login')
}

onMounted(() => {
  void loadDashboard()
})
</script>

<template>
  <div class="admin-dashboard-page">
    <div class="admin-dashboard-shell">
      <section class="hero">
        <div>
          <p class="eyebrow">Operations Console</p>
          <h1>{{ welcomeTitle }}</h1>
          <p>
            统一管理题目、题单、用户和系统设置，先把登录链路和基础运营入口稳定下来。
          </p>
        </div>
        <n-button secondary strong type="primary" @click="handleLogout">退出登录</n-button>
      </section>

      <n-result
        v-if="errorMessage"
        status="error"
        title="加载失败"
        :description="errorMessage"
      />

      <template v-else-if="pending">
        <n-space vertical :size="18">
          <n-skeleton text :repeat="2" />
          <n-grid :cols="2" :x-gap="18" :y-gap="18">
            <n-grid-item v-for="item in 4" :key="item">
              <n-card :bordered="false">
                <n-skeleton text :repeat="4" />
              </n-card>
            </n-grid-item>
          </n-grid>
        </n-space>
      </template>

      <template v-else-if="dashboard">
        <n-grid class="content-grid" :cols="24" :x-gap="20" :y-gap="20">
          <n-grid-item :span="10">
            <n-card class="panel" :bordered="false" title="当前管理员">
              <n-space vertical :size="14">
                <div class="identity-name">{{ dashboard.admin_user.display_name }}</div>
                <div class="identity-meta">{{ dashboard.admin_user.email }}</div>
                <n-tag type="success" size="small">{{ dashboard.admin_user.status }}</n-tag>
              </n-space>
            </n-card>
          </n-grid-item>

          <n-grid-item :span="14">
            <n-card class="panel" :bordered="false" title="快捷入口">
              <div class="quick-links">
                <div v-for="link in dashboard.quick_links" :key="link.path" class="quick-link-card">
                  <div class="quick-link-title">{{ link.label }}</div>
                  <div class="quick-link-path">{{ link.path }}</div>
                </div>
              </div>
            </n-card>
          </n-grid-item>
        </n-grid>
      </template>
    </div>
  </div>
</template>

<style scoped>
.admin-dashboard-page {
  min-height: 100vh;
  padding: 28px;
  background:
    linear-gradient(180deg, rgba(16, 41, 66, 0.95) 0%, rgba(23, 57, 84, 0.92) 28%, #eef2f5 28%, #eef2f5 100%);
}

.admin-dashboard-shell {
  max-width: 1200px;
  margin: 0 auto;
}

.hero {
  display: flex;
  justify-content: space-between;
  gap: 20px;
  align-items: flex-start;
  padding: 36px;
  margin-bottom: 24px;
  border-radius: 28px;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.14), rgba(255, 255, 255, 0.06));
  color: #f7fafc;
  backdrop-filter: blur(10px);
}

.hero h1 {
  margin: 8px 0 12px;
  font-size: 36px;
}

.hero p {
  max-width: 640px;
  margin: 0;
  line-height: 1.7;
}

.eyebrow {
  margin: 0;
  letter-spacing: 0.2em;
  text-transform: uppercase;
  font-size: 12px;
  color: #ffd27f;
}

.content-grid {
  align-items: stretch;
}

.panel {
  border-radius: 24px;
  box-shadow: 0 20px 50px rgba(15, 36, 58, 0.08);
}

.identity-name {
  font-size: 28px;
  font-weight: 700;
  color: #1f3142;
}

.identity-meta {
  color: #607181;
}

.quick-links {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 14px;
}

.quick-link-card {
  padding: 18px;
  border-radius: 18px;
  background: linear-gradient(160deg, #f6efe3, #ebf2f7);
}

.quick-link-title {
  font-size: 18px;
  font-weight: 700;
  color: #1f3142;
}

.quick-link-path {
  margin-top: 8px;
  color: #647485;
}

@media (max-width: 768px) {
  .admin-dashboard-page {
    padding: 16px;
  }

  .hero {
    flex-direction: column;
    padding: 24px;
  }
}
</style>
