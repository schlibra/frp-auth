<script setup lang="ts">
import { h, onMounted, type Ref, ref } from 'vue'
import { type DataTableColumn, NTag } from 'naive-ui'
import axios from 'axios'
import { getToken } from '@/utils/token.ts'
import { useDialog } from 'naive-ui'
import type Response from '@/model/response.ts'
import { dialogError } from '@/utils/dialog.ts'
import { calcSize } from '@/utils/size.ts'
import { useProxiesStore } from '@/stores'

const dialog = useDialog()
const proxies = useProxiesStore()
const columns: Ref<DataTableColumn[]> = ref([
  {
    title: '状态',
    key: 'status.phase',
    render(row) {
      const _status = row.status.phase
      return h(
        NTag,
        {
          type: _status == 'online' ? 'success' : _status == 'offline' ? 'default' : 'warning',
          round: true,
        },
        _status == 'online' ? '在线' : _status == 'offline' ? '离线' : '未知',
      )
    },
  },
  {
    title: '名称',
    key: 'name',
    width: 200,
    render(row) {
      let full_name = row.name
      const user = row.user
      full_name = full_name.replace(`${user}.`, '')
      return h('span', {}, full_name)
    },
  },
  {
    title: '远程端口',
    key: 'remotePort',
    width: 100,
    render(row) {
      const type = row.spec.type
      const port = row.spec[type].remotePort
      return h('spane', {}, port)
    },
  },
  {
    title: 'Token',
    key: 'user',
    width: 100,
  },
  {
    title: '类型',
    key: 'spec.type',
    width: 70,
    render(row) {
      return h(
        NTag,
        {
          type: 'info',
          round: true,
        },
        row.spec.type,
      )
    },
  },
  {
    title: '今日入站流量',
    key: 'status.todayTrafficIn',
    width: 110,
    render(row) {
      return h('span', {}, calcSize(row.status.todayTrafficIn))
    },
  },
  {
    title: '今日出站流量',
    key: 'status.todayTrafficOut',
    width: 110,
    render(row) {
      return h('span', {}, calcSize(row.status.todayTrafficOut))
    },
  },
  {
    title: '当前连接数量',
    key: 'status.curConns',
    width: 110,
  },
  {
    title: '上次连接时间',
    key: 'status.lastStartAt',
    width: 150,
    render(row) {
      const time = new Date()
      time.setTime(row.status.lastStartAt * 1000)
      return h('span', {}, time.toLocaleString())
    },
  },
  {
    title: '上次断开时间',
    key: 'status.lastCloseAt',
    width: 150,
    render(row) {
      const _t = row.status.lastCloseAt
      if (_t) {
        const time = new Date()
        time.setTime(_t * 1000)
        return h('span', {}, time.toLocaleString())
      } else {
        return h('span', {}, '-')
      }
    },
  },
])

async function loadProxy() {
  const res = await axios.get('/api/port', {
    headers: {
      Authorization: `Bearer ${getToken()}`,
    },
  })
  const data: Response = res.data
  if (data.status === 200) {
    proxies.proxies = data.data.proxies
  } else {
    dialogError(dialog, '数据获取失败', data.msg)
  }
}

onMounted(async () => {
  await loadProxy()
})
</script>

<template>
  <n-card>
    <template #header>
      <h3>映射管理</h3>
    </template>
    <n-space vertical>
      <n-button type="primary" size="large" @click="loadProxy()">刷新数据</n-button>
      <n-data-table :columns="columns" :data="proxies.proxies"></n-data-table>
    </n-space>
  </n-card>
</template>

<style scoped></style>
