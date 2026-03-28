<script setup lang="ts">
import { NCard, NSpace, NStatistic } from 'naive-ui'
</script>

<template>
  <div class="page-section route-page">
    <h1 class="page-title">我的进度</h1>
    <p class="page-subtitle">本页先保留进度占位，为后续用户题目统计接入留出入口。</p>
    <n-card class="stack-gap">
      <n-space justify="space-around">
        <n-statistic label="已通过" :value="18" />
        <n-statistic label="尝试中" :value="7" />
        <n-statistic label="本周提交" :value="12" />
      </n-space>
    </n-card>
  </div>
</template>

<style scoped>
.route-page {
  padding: 32px 0 48px;
}

.stack-gap {
  margin-top: 20px;
}
</style>
