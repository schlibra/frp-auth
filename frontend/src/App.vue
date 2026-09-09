<script setup lang="ts">
import { darkTheme, zhCN } from 'naive-ui'
</script>

<template>
  <n-config-provider :theme="darkTheme" :locale="zhCN">
    <n-dialog-provider>
      <n-modal-provider>
        <router-view class="main"></router-view>
      </n-modal-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>

<style scoped>
.main {
  height: 100vh;
}
</style>
