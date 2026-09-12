<script setup lang="ts">
import { ref } from 'vue'

const pathTable = {
  '/': '首页',
  '/token': 'Token管理',
  '/port': '端口规则管理',
  '/client': '客户端管理',
  '/proxy': '映射管理',
  '/config': '配置生成',
  '/user': '个人设置',
  '/admin/token': '管理员 - Token管理',
  '/admin/port': '管理员 - 端口规则管理',
  '/admin/client': '管理员 - 客户端管理',
  '/admin/proxy': '管理员 - 映射管理',
  '/admin/user': '管理员 - 用户管理',
}
const path = ref(location.pathname)
setInterval(() => {
  path.value = location.pathname
  document.title = pathTable[path.value]
}, 100)
</script>

<template>
  <n-page-header class="header">
    <template #title>
      <h3>Frp-Auth</h3>
    </template>
    <template #subtitle>
      <span>{{ pathTable[path] }}</span>
    </template>
  </n-page-header>
</template>

<style scoped>
.header {
  padding-left: 16px;
}
</style>
