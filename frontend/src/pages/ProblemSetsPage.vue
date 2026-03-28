<script setup lang="ts">
import { NCard, NGrid, NGridItem, NTag } from 'naive-ui'

import { useProblemSetsQuery } from '../features/problem-sets/queries'

const { data } = useProblemSetsQuery()
</script>

<template>
  <div class="page-section route-page">
    <h1 class="page-title">题单</h1>
    <p class="page-subtitle">先从一组固定专题进入刷题节奏。</p>
    <n-grid cols="1 s:2" responsive="screen" :x-gap="16" :y-gap="16">
      <n-grid-item v-for="item in data ?? []" :key="item.slug">
        <RouterLink :to="`/problem-sets/${item.slug}`">
          <n-card hoverable>
            <n-tag>{{ item.problemCount }} 题</n-tag>
            <h3>{{ item.name }}</h3>
            <p>{{ item.description }}</p>
          </n-card>
        </RouterLink>
      </n-grid-item>
    </n-grid>
  </div>
</template>

<style scoped>
.route-page {
  padding: 32px 0 48px;
}
</style>
