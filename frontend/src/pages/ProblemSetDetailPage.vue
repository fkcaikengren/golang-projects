<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { NCard, NList, NListItem, NSpace, NTag } from 'naive-ui'

import { useProblemSetDetailQuery } from '../features/problem-sets/queries'
import { useProblemsQuery } from '../features/problems/queries'

const route = useRoute()
const slug = computed(() => String(route.params.slug ?? 'hot-100'))
const { data: detail } = useProblemSetDetailQuery(slug.value)
const { data: problems } = useProblemsQuery()
</script>

<template>
  <div class="page-section route-page">
    <n-card>
      <n-space vertical>
        <n-tag type="info">{{ detail?.problemCount ?? 0 }} 题</n-tag>
        <h1 class="page-title">{{ detail?.name }}</h1>
        <p class="page-subtitle">{{ detail?.description }}</p>
      </n-space>
    </n-card>
    <n-card title="题目列表" class="stack-gap">
      <n-list>
        <n-list-item v-for="item in problems ?? []" :key="item.slug">
          <div class="problem-line">
            <RouterLink :to="`/problems/${item.slug}`">{{ item.title }}</RouterLink>
            <n-space>
              <n-tag>{{ item.difficulty }}</n-tag>
              <n-tag type="success">{{ item.status }}</n-tag>
            </n-space>
          </div>
        </n-list-item>
      </n-list>
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

.problem-line {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}
</style>
