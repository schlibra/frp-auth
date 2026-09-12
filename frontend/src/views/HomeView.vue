<script setup lang="ts">
import {
  KeyOutline as TokenIcon,
  GlobeOutline as PortIcon,
  DesktopOutline as ClientIcon,
  RepeatOutline as ProxyIcon,
} from '@vicons/ionicons5'
import { onMounted, ref } from 'vue'
import { getToken } from '@/utils/token.ts'
import axios from 'axios'
import type Response from '@/model/response.ts'
import { dialogError } from '@/utils/dialog.ts'
import { useDialog } from 'naive-ui'
import router from '@/router'
import { useTokensStore, usePortsStore, useClientsStore } from '@/stores'
import { useProxiesStore } from '@/stores/proxies.ts'

const dialog = useDialog()
const tokens = useTokensStore()
const ports = usePortsStore()
const clients = useClientsStore()
const proxies = useProxiesStore()

onMounted(async () => {
  let res = await axios.get('/api/token', {
    headers: {
      Authorization: 'Bearer ' + getToken(),
    },
  })
  let data: Response = res.data
  if (data.status == 200) {
    tokens.tokens = data.data.tokens
  } else if (data.status == 401) {
    router.push('/login')
  } else {
    dialogError(dialog, '数据获取失败', data.msg)
  }
  res = await axios.get('/api/rule', {
    headers: {
      Authorization: 'Bearer ' + getToken(),
    },
  })
  data = res.data
  if (data.status == 200) {
    ports.ports = data.data.ports
  } else if (data.status == 401) {
    router.push('/login')
  } else {
    dialogError(dialog, '数据获取失败', data.msg)
  }
  res = await axios.get('/api/port', {
    headers: {
      Authorization: 'Bearer ' + getToken(),
    },
  })
  data = res.data
  if (data.status == 200) {
    clients.clients = data.data.clients
  } else if (data.status == 401) {
    router.push('/login')
  } else {
    dialogError(dialog, '数据获取失败', data.msg)
  }
  res = await axios.get('/api/port', {
    headers: {
      Authorization: 'Bearer ' + getToken(),
    },
  })
  data = res.data
  if (data.status == 200) {
    proxies.proxies = data.data.proxies
  } else if (data.status == 401) {
    router.push('/login')
  } else {
    dialogError(dialog, '数据获取失败', data.msg)
  }
})
</script>

<template>
  <n-card>
    <template #header>
      <h3>首页</h3>
    </template>
    <n-flex>
      <n-statistic label="Token数量" :value="tokens.count">
        <template #prefix>
          <n-icon>
            <TokenIcon />
          </n-icon>
        </template>
      </n-statistic>
      <n-statistic label="端口规则数量" :value="ports.count">
        <template #prefix>
          <n-icon>
            <PortIcon />
          </n-icon>
        </template>
      </n-statistic>
      <n-statistic label="客户端数量" :value="clients.count">
        <template #prefix>
          <n-icon>
            <ClientIcon />
          </n-icon>
        </template>
      </n-statistic>
      <n-statistic label="映射数量" :value="proxies.count">
        <template #prefix>
          <n-icon>
            <ProxyIcon />
          </n-icon>
        </template>
      </n-statistic>
    </n-flex>
  </n-card>
</template>

<style scoped></style>
