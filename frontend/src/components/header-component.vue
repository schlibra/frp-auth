<script setup lang="ts">
const props = defineProps(['title'])
</script>

<template>
  <n-page-header class="header">
    <template #title>
      <h3>Frp-Auth</h3>
    </template>
    <template #subtitle>
      <span>{{props.title}}</span>
    </template>
  </n-page-header>
</template>

<style scoped>
.header {
  padding-left: 16px;
}
</style>
