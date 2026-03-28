<script setup lang="ts">
import { NButton, NCard, NSpace, NText } from 'naive-ui'
</script>

<template>
  <n-card class="hero-panel" :bordered="false">
    <n-space vertical :size="18">
      <n-text depth="3">Go OJ MVP</n-text>
      <h1 class="hero-title">开始刷题</h1>
      <p class="hero-copy">
        面向中文开发者的在线判题起步站点。先进入题单和题库，再逐步接入真实判题流程。
      </p>
      <n-space>
        <RouterLink to="/problem-sets">
          <n-button type="primary" size="large">进入热门题单</n-button>
        </RouterLink>
        <RouterLink to="/problems">
          <n-button size="large" tertiary>浏览题库</n-button>
        </RouterLink>
      </n-space>
    </n-space>
  </n-card>
</template>

<style scoped>
.hero-panel {
  background: linear-gradient(135deg, #0f172a 0%, #1d4ed8 100%);
  color: #fff;
}

.hero-title {
  margin: 0;
  font-size: 56px;
  line-height: 1.05;
}

.hero-copy {
  margin: 0;
  max-width: 560px;
  color: rgba(255, 255, 255, 0.82);
}
</style>
