<script setup lang="ts">
import { NCard, NGrid, NGridItem, NTag } from 'naive-ui'

import { useProblemSetsQuery } from '../../features/problem-sets/queries'

const { data } = useProblemSetsQuery()
</script>

<template>
  <section>
    <div class="section-head">
      <h2>热门题单</h2>
      <p>从固定专题切入，先把主链路跑起来。</p>
    </div>
    <n-grid cols="1 s:2" responsive="screen" :x-gap="16" :y-gap="16">
      <n-grid-item v-for="item in data ?? []" :key="item.slug">
        <RouterLink :to="`/problem-sets/${item.slug}`">
          <n-card hoverable>
            <n-tag type="info" size="small">{{ item.problemCount }} 题</n-tag>
            <h3>{{ item.name }}</h3>
            <p>{{ item.description }}</p>
          </n-card>
        </RouterLink>
      </n-grid-item>
    </n-grid>
  </section>
</template>

<style scoped>
.section-head {
  margin-bottom: 16px;
}

h2,
h3,
p {
  margin-top: 0;
}
</style>
