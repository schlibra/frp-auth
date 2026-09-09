<script setup lang="ts">
import router from '@/router'

function goHome() {
  router.push('/')
}
</script>

<template>
  <n-layout>
    <n-layout-content>
      <n-card class="main">
        <n-result status="404" title="404 页面不存在">
          <template #default>
            <n-button size="large" @click="goHome">返回首页</n-button>
          </template>
        </n-result>
      </n-card>
    </n-layout-content>
  </n-layout>
</template>

<style scoped>
.main {
  text-align: center;
  padding: 16px;
}
</style>
