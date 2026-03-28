<script setup lang="ts">
import { computed, h, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { NCard, NDataTable, NInput, NSelect, NSpace, NTag } from 'naive-ui'

import { useProblemsQuery } from '../features/problems/queries'

const difficulty = ref<string | null>(null)
const keyword = ref('')
const { data } = useProblemsQuery()

const rows = computed(() =>
  (data.value ?? []).filter((item) => {
    const difficultyMatch = !difficulty.value || item.difficulty === difficulty.value
    const keywordMatch = !keyword.value || item.title.includes(keyword.value)
    return difficultyMatch && keywordMatch
  }),
)

const columns = [
  {
    title: '题目',
    key: 'title',
    render: (row: { slug: string; title: string }) =>
      h(
        RouterLink,
        { to: `/problems/${row.slug}` },
        { default: () => row.title },
      ),
  },
  { title: '难度', key: 'difficulty' },
  {
    title: '标签',
    key: 'tags',
    render: (row: { tags: string[] }) => row.tags.join(' / '),
  },
  { title: '来源', key: 'source' },
  { title: '状态', key: 'status' },
]
</script>

<template>
  <div class="page-section route-page">
    <h1 class="page-title">题库</h1>
    <p class="page-subtitle">按关键词和难度做最小筛选，先跑通浏览主流程。</p>
    <n-card class="stack-gap">
      <n-space>
        <n-input v-model:value="keyword" placeholder="搜索题目" />
        <n-select
          v-model:value="difficulty"
          clearable
          placeholder="按难度筛选"
          :options="[
            { label: 'easy', value: 'easy' },
            { label: 'medium', value: 'medium' },
            { label: 'hard', value: 'hard' },
          ]"
        />
      </n-space>
    </n-card>
    <n-card class="stack-gap" title="题目列表">
      <n-data-table :columns="columns" :data="rows" />
      <div class="quick-links">
        <RouterLink v-for="item in rows" :key="item.slug" :to="`/problems/${item.slug}`">
          <n-tag round>{{ item.title }}</n-tag>
        </RouterLink>
      </div>
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

.quick-links {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 16px;
}
</style>
