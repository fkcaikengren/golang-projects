<script setup lang="ts">
import { NGrid, NGridItem, NSpace } from 'naive-ui'

import FeaturedSets from '../widgets/home/FeaturedSets.vue'
import HeroPanel from '../widgets/home/HeroPanel.vue'
import RecentProblems from '../widgets/home/RecentProblems.vue'
</script>

<template>
  <div class="page-section home-page">
    <n-space vertical :size="24">
      <HeroPanel />
      <n-grid cols="1 l:3" responsive="screen" :x-gap="20" :y-gap="20">
        <n-grid-item span="2">
          <FeaturedSets />
        </n-grid-item>
        <n-grid-item>
          <RecentProblems />
        </n-grid-item>
      </n-grid>
    </n-space>
  </div>
</template>

<style scoped>
.home-page {
  padding: 32px 0 48px;
}
</style>
