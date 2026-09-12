<script setup lang="ts">
import HeaderComponent from '@/components/header-component.vue'
import MenuComponent from '@/components/menu-component.vue'
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

const dialog = useDialog()
const tokenCount = ref(0)
const portCount = ref(0)
const clientCount = ref(0)
const proxyCount = ref(0)

onMounted(async () => {
  let res = await axios.get('/api/token', {
    headers: {
      Authorization: 'Bearer ' + getToken(),
    },
  })
  let data: Response = res.data
  if (data.status == 200) {
    tokenCount.value = data.data.tokens.length
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
    portCount.value = data.data.ports.length
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
    clientCount.value = data.data.clients.length
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
    proxyCount.value = data.data.proxies.length
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
      <n-statistic label="Token数量" :value="tokenCount">
        <template #prefix>
          <n-icon>
            <TokenIcon />
          </n-icon>
        </template>
      </n-statistic>
      <n-statistic label="端口规则数量" :value="portCount">
        <template #prefix>
          <n-icon>
            <PortIcon />
          </n-icon>
        </template>
      </n-statistic>
      <n-statistic label="客户端数量" :value="clientCount">
        <template #prefix>
          <n-icon>
            <ClientIcon />
          </n-icon>
        </template>
      </n-statistic>
      <n-statistic label="映射数量" :value="proxyCount">
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
