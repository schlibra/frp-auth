<script setup lang="ts">
import HeaderComponent from '@/components/header-component.vue'
import MenuComponent from '@/components/menu-component.vue'
import { computed, ref } from 'vue'

const preConfig = ref('')
const proxyName = ref('')
const proxyType = ref('')
const proxyTypeOptions = ref([
  {
    label: 'TCP',
    value: 'tcp',
  },
  {
    label: 'UDP',
    value: 'udp',
  }
])
const localIP = ref('')
const localPort = ref('')
const remotePort = ref('')
const resultRef = ref(null)
const result = computed(
  () => proxyName.value && proxyType.value && localIP.value && localPort.value && remotePort.value ? `${preConfig.value}

[[proxies]]
name = "${proxyName.value}"
type = "${proxyType.value}"
localIP = "${localIP.value}"
localPort = "${localPort.value}"
remotePort = "${remotePort.value}"
` : '',
)
function copy() {
  resultRef.value.select()
  document.execCommand('copy')
}
function clear() {

}
</script>

<template>
  <n-layout>
    <n-layout-header>
      <header-component title="配置生成"></header-component>
    </n-layout-header>
    <n-layout has-sider>

        <menu-component></menu-component>

      <n-layout-content>
        <n-card>
          <template #header>
            <h3>配置生成</h3>
          </template>
          <n-form>
            <n-form-item label="前置配置">
              <n-input v-model:value="preConfig" type="textarea" rows="5"></n-input>
            </n-form-item>
            <n-flex>
              <n-form-item label="映射名称">
                <n-input v-model:value="proxyName"></n-input>
              </n-form-item>
              <n-form-item label="映射类型" style="min-width: 150px">
                <n-select :options="proxyTypeOptions" v-model:value="proxyType"></n-select>
              </n-form-item>
              <n-form-item label="本地IP">
                <n-input v-model:value="localIP"></n-input>
              </n-form-item>
              <n-form-item label="本地端口">
                <n-input-number v-model:value="localPort"></n-input-number>
              </n-form-item>
              <n-form-item label="远程端口">
                <n-input-number v-model:value="remotePort"></n-input-number>
              </n-form-item>
            </n-flex>
            <n-form-item label="生成配置">
              <n-input ref="resultRef" type="textarea" rows="8" :value="result"></n-input>
            </n-form-item>
          </n-form>
          <template #footer>
            <n-space justify="center">
              <n-button type="primary" size="large" @click="copy()">复制</n-button>
              <n-button type="info" size="large" @click="clear()">清空</n-button>
            </n-space>
          </template>
        </n-card>
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<style scoped></style>
