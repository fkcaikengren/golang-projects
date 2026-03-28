<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { AxiosError } from 'axios'
import { NAlert, NButton, NCard, NForm, NFormItem, NInput, NSpace } from 'naive-ui'

import { adminLogin } from '../features/admin-auth/api'
import { useAdminAuthStore } from '../features/admin-auth/store'

const router = useRouter()
const adminAuthStore = useAdminAuthStore()

const email = ref('admin@go-oj.dev')
const password = ref('Admin@123456')
const pending = ref(false)
const errorMessage = ref('')

const submitLabel = computed(() => (pending.value ? '登录中...' : '进入管理后台'))

async function handleSubmit() {
  pending.value = true
  errorMessage.value = ''

  try {
    const result = await adminLogin({
      email: email.value,
      password: password.value,
    })

    adminAuthStore.setSession(result)
    await router.push('/admin')
  } catch (error) {
    if (error instanceof AxiosError) {
      errorMessage.value = error.response?.data?.message ?? '管理登录失败，请检查账号密码。'
    } else {
      errorMessage.value = '管理登录失败，请稍后重试。'
    }
  } finally {
    pending.value = false
  }
}
</script>

<template>
  <div class="admin-login-page">
    <div class="admin-login-panel">
      <n-card class="admin-login-card" :bordered="false">
        <n-space vertical :size="20">
          <div>
            <p class="eyebrow">Go OJ Admin</p>
            <h1 class="title">管理后台登录</h1>
            <p class="subtitle">使用默认超级管理员账号进入控制台，开始维护题目、题单和系统配置。</p>
          </div>

          <n-alert v-if="errorMessage" type="error" :show-icon="false">
            {{ errorMessage }}
          </n-alert>

          <n-form label-placement="top">
            <n-form-item label="管理员邮箱">
              <n-input v-model:value="email" placeholder="输入管理员邮箱" />
            </n-form-item>
            <n-form-item label="密码">
              <n-input
                v-model:value="password"
                type="password"
                show-password-on="click"
                placeholder="输入密码"
                @keyup.enter="handleSubmit"
              />
            </n-form-item>
            <n-button type="primary" block :loading="pending" @click="handleSubmit">
              {{ submitLabel }}
            </n-button>
          </n-form>
        </n-space>
      </n-card>

      <div class="admin-login-tip">
        <span>默认账号：</span>
        <strong>{{ email }}</strong>
        <span> / </span>
        <strong>{{ password }}</strong>
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin-login-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 32px 20px;
  background:
    radial-gradient(circle at top left, rgba(244, 181, 98, 0.3), transparent 30%),
    radial-gradient(circle at bottom right, rgba(31, 78, 121, 0.32), transparent 35%),
    linear-gradient(145deg, #f5f1e8 0%, #e5edf3 50%, #d8e1e6 100%);
}

.admin-login-panel {
  width: min(460px, 100%);
}

.admin-login-card {
  border-radius: 28px;
  box-shadow: 0 24px 80px rgba(34, 52, 70, 0.14);
}

.eyebrow {
  margin: 0 0 8px;
  font-size: 12px;
  letter-spacing: 0.24em;
  text-transform: uppercase;
  color: #9a5d18;
}

.title {
  margin: 0;
  font-size: 32px;
  color: #1f3142;
}

.subtitle {
  margin: 12px 0 0;
  line-height: 1.7;
  color: #536273;
}

.admin-login-tip {
  margin-top: 16px;
  padding: 14px 16px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.72);
  color: #405061;
  text-align: center;
}
</style>
