<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NButton, NCard, NForm, NFormItem, NInput, NSpace } from 'naive-ui'

import { login, register } from '../features/auth/api'
import { useAuthStore } from '../features/auth/store'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const email = ref('demo@go-oj.dev')
const password = ref('password123')
const nickname = ref('DemoUser')
const mode = computed(() => (route.path === '/register' ? 'register' : 'login'))

async function handleSubmit() {
  const payload = {
    email: email.value,
    password: password.value,
    nickname: nickname.value,
  }

  const result =
    mode.value === 'register' ? await register(payload) : await login(payload)

  authStore.setSession(result)
  await router.push('/problem-sets')
}
</script>

<template>
  <div class="page-section auth-page">
    <n-card class="auth-card" :bordered="false">
      <n-space vertical :size="20">
        <div>
          <h1 class="page-title">登录 / 注册</h1>
          <p class="page-subtitle">当前先使用 mock 流程打通身份页和登录态。</p>
        </div>
        <n-form label-placement="top">
          <n-form-item label="邮箱">
            <n-input v-model:value="email" placeholder="输入邮箱" />
          </n-form-item>
          <n-form-item v-if="mode === 'register'" label="昵称">
            <n-input v-model:value="nickname" placeholder="输入昵称" />
          </n-form-item>
          <n-form-item label="密码">
            <n-input v-model:value="password" type="password" placeholder="输入密码" />
          </n-form-item>
          <n-button type="primary" block @click="handleSubmit">
            {{ mode === 'register' ? '创建账号' : '登录并开始刷题' }}
          </n-button>
        </n-form>
      </n-space>
    </n-card>
  </div>
</template>

<style scoped>
.auth-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 72px);
}

.auth-card {
  width: min(480px, 100%);
}
</style>
