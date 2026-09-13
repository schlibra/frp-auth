<script setup lang="ts">
import { ref } from 'vue'
import { useConfigStore } from '@/stores'

const config = useConfigStore()
const proxyTypeOptions = ref([
  {
    label: 'TCP',
    value: 'tcp',
  },
  {
    label: 'UDP',
    value: 'udp',
  },
])
const resultRef = ref(null)

function copy() {
  resultRef.value.select()
  document.execCommand('copy')
}
</script>

<template>
  <n-card>
    <template #header>
      <h3>配置生成</h3>
    </template>
    <n-form>
      <n-form-item label="前置配置">
        <n-input v-model:value="config.preConfig" type="textarea" rows="4"></n-input>
      </n-form-item>
      <n-flex>
        <n-form-item label="映射名称">
          <n-input v-model:value="config.proxyName"></n-input>
        </n-form-item>
        <n-form-item label="映射类型" style="min-width: 150px">
          <n-select :options="proxyTypeOptions" v-model:value="config.proxyType"></n-select>
        </n-form-item>
        <n-form-item label="本地IP">
          <n-input v-model:value="config.localIP"></n-input>
        </n-form-item>
        <n-form-item label="本地端口">
          <n-input-number v-model:value="config.localPort"></n-input-number>
        </n-form-item>
        <n-form-item label="远程端口">
          <n-input-number v-model:value="config.remotePort"></n-input-number>
        </n-form-item>
      </n-flex>
      <n-form-item label="生成配置">
        <n-input ref="resultRef" type="textarea" rows="8" :value="config.result"></n-input>
      </n-form-item>
    </n-form>
    <template #footer>
      <n-space justify="center">
        <n-button type="success" size="large" @click="config.next()">编辑下一条</n-button>
        <n-button type="info" size="large" @click="copy()">复制</n-button>
        <n-button type="warning" size="large" @click="config.clear()">清空</n-button>
      </n-space>
    </template>
  </n-card>
</template>

<style scoped></style>
