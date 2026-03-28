<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import { NButton, NCard, NCode, NInput, NSelect, NSpace, NTag } from 'naive-ui'

import { useProblemDetailQuery, useSubmitProblemMutation } from '../features/problems/queries'

const route = useRoute()
const slug = computed(() => String(route.params.slug ?? 'two-sum'))
const language = ref('golang')
const code = ref('package main\n\nfunc main() {\n    // TODO\n}\n')
const { data } = useProblemDetailQuery(slug.value)
const submitMutation = useSubmitProblemMutation()

async function handleSubmit() {
  await submitMutation.mutateAsync({
    slug: slug.value,
    code: code.value,
    language: language.value,
  })
}
</script>

<template>
  <div class="page-section route-page">
    <div class="detail-grid">
      <n-card>
        <n-space vertical :size="16">
          <div>
            <h1 class="page-title">{{ data?.title }}</h1>
            <n-space>
              <n-tag>{{ data?.difficulty }}</n-tag>
              <n-tag v-for="tag in data?.tags ?? []" :key="tag" type="info" size="small">
                {{ tag }}
              </n-tag>
            </n-space>
          </div>
          <p>{{ data?.description }}</p>
          <section>
            <h3>输入说明</h3>
            <p>{{ data?.inputDescription }}</p>
          </section>
          <section>
            <h3>输出说明</h3>
            <p>{{ data?.outputDescription }}</p>
          </section>
          <section>
            <h3>示例</h3>
            <n-code :code="data?.sampleInput ?? ''" language="text" />
            <n-code :code="data?.sampleOutput ?? ''" language="text" />
          </section>
          <section>
            <h3>提示</h3>
            <p>{{ data?.hint }}</p>
          </section>
        </n-space>
      </n-card>
      <n-card title="代码编辑与提交">
        <n-space vertical :size="16">
          <n-select
            v-model:value="language"
            :options="[
              { label: 'Go', value: 'golang' },
              { label: 'JavaScript', value: 'javascript' },
            ]"
          />
          <n-input
            v-model:value="code"
            type="textarea"
            :autosize="{ minRows: 18, maxRows: 24 }"
          />
          <n-button type="primary" @click="handleSubmit">提交代码</n-button>
          <p v-if="submitMutation.data.value">{{ submitMutation.data.value.message }}</p>
        </n-space>
      </n-card>
    </div>
  </div>
</template>

<style scoped>
.route-page {
  padding: 32px 0 48px;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1.35fr 1fr;
  gap: 20px;
}

@media (max-width: 900px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }
}
</style>
