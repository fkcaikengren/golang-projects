<script setup lang="ts">
import { NCard, NList, NListItem, NTag, NSpace } from 'naive-ui'

import { useProblemsQuery } from '../../features/problems/queries'

const { data } = useProblemsQuery()
</script>

<template>
  <n-card title="最近题目">
    <n-list>
      <n-list-item v-for="item in data ?? []" :key="item.slug">
        <div class="problem-line">
          <RouterLink :to="`/problems/${item.slug}`">{{ item.title }}</RouterLink>
          <n-space>
            <n-tag size="small">{{ item.difficulty }}</n-tag>
            <n-tag size="small" type="success">{{ item.status }}</n-tag>
          </n-space>
        </div>
      </n-list-item>
    </n-list>
  </n-card>
</template>

<style scoped>
.problem-line {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}
</style>
