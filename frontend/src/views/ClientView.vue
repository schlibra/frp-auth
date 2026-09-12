<script setup lang="ts">
import { h, onMounted, type Ref, ref } from 'vue'
import type Response from '@/model/response.ts'
import axios from 'axios'
import { getToken } from '@/utils/token.ts'
import { type DataTableColumn, NTag, useDialog } from 'naive-ui'
import { dialogError } from '@/utils/dialog.ts'
import { useClientsStore } from '@/stores'

const dialog = useDialog()
const clients = useClientsStore()
const columns: Ref<DataTableColumn[]> = ref([
  {
    title: '状态',
    key: 'online',
    width: 70,
    render(row) {
      return h(
        NTag,
        {
          size: 'large',
          type: row.online ? 'success' : 'default',
          round: true,
        },
        row.online ? '在线' : '离线',
      )
    },
  },
  {
    title: 'Token',
    key: 'user',
    width: 100
  },
  {
    title: '客户端IP',
    key: 'clientIP',
    width: 150
  },
  {
    title: '主机名',
    key: 'hostname',
    width: 150
  },
  {
    title: '客户端ID',
    key: 'clientID',
    width: 160
  },
  {
    title: '客户端版本',
    key: 'version',
    width: 100
  },
  {
    title: '首次连接时间',
    key: 'firstConnectedAt',
    width: 150,
    render(row) {
      const time = new Date()
      time.setTime(row.firstConnectedAt * 1000)
      return h('span', {}, time.toLocaleString())
    },
  },
  {
    title: '上次连接时间',
    key: 'lastConnectedAt',
    width: 150,
    render(row) {
      const time = new Date()
      time.setTime(row.lastConnectedAt * 1000)
      return h('span', {}, time.toLocaleString())
    },
  },
])

async function loadClient() {
  const res = await axios.get('/api/port', {
    headers: {
      Authorization: `Bearer ${getToken()}`,
    },
  })
  const data: Response = res.data
  if (data.status === 200) {
    clients.clients = data.data.clients
  } else {
    dialogError(dialog, '数据获取失败', data.msg)
  }
}

onMounted(async () => {
  await loadClient()
})
</script>

<template>
  <n-card>
    <template #header>
      <h3>客户端管理</h3>
    </template>
    <n-space vertical>
      <n-button type="primary" size="large" @click="loadClient()">刷新数据</n-button>
      <n-data-table :columns="columns" :data="clients.clients"></n-data-table>
    </n-space>
  </n-card>
</template>

<style scoped></style>
