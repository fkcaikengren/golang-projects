<script setup lang="ts">
import { NCard, NDataTable } from 'naive-ui'

import { useSubmissionsQuery } from '../features/submissions/queries'

const { data } = useSubmissionsQuery()

const columns = [
  { title: '题目', key: 'problemTitle' },
  { title: '语言', key: 'language' },
  { title: '结果', key: 'result' },
  { title: '提交时间', key: 'submittedAt' },
]
</script>

<template>
  <div class="page-section route-page">
    <h1 class="page-title">提交记录</h1>
    <p class="page-subtitle">当前展示 mock 驱动的个人提交历史。</p>
    <n-card class="stack-gap">
      <n-data-table :columns="columns" :data="data ?? []" />
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
