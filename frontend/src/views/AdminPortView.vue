<script setup lang="ts">
import { h, onMounted, type Ref, ref } from 'vue'
import axios from 'axios'
import { getToken } from '@/utils/token.ts'
import type Response from '@/model/response.ts'
import { type DataTableColumn, NButton, NSpace, useDialog } from 'naive-ui'
import { dialogError } from '@/utils/dialog.ts'
import router from '@/router'
import { useAdminPortsStore } from '@/stores'

const dialog = useDialog()
const ports = useAdminPortsStore()
const modalShow = ref(false)
const createPort = ref(false)
const modalPortId = ref(0)
const modalPortMin = ref(0)
const modalPortMax = ref(0)
const modalPortToken = ref('')
const modalPortUser = ref('')
const columns: Ref<DataTableColumn[]> = ref([
  {
    title: 'ID',
    key: 'id',
    width: 70
  },
  {
    title: '最小端口',
    key: 'min',
    width: 80
  },
  {
    title: '最大端口',
    key: 'max',
    width: 80
  },
  {
    title: 'Token',
    key: 'token',
    width: 120
  },
  {
    title: '用户',
    key: 'user',
    width: 80
  },
  {
    title: '操作',
    key: 'action',
    width: 150,
    render(row) {
      return h(NSpace, {}, [
        h(
          NButton,
          {
            type: 'warning',
            onClick() {
              openUpdateModal(row)
            },
          },
          '编辑',
        ),
        h(
          NButton,
          {
            type: 'error',
            onClick() {
              deletePort(row)
            },
          },
          '删除',
        ),
      ])
    },
  },
])

function openCreateModal() {
  modalPortId.value = 0
  modalPortMin.value = 0
  modalPortMax.value = 0
  modalPortToken.value = ''
  modalPortUser.value = ''
  createPort.value = true
  modalShow.value = true
}

function openUpdateModal(row) {
  modalPortId.value = row.id
  modalPortMin.value = row.min
  modalPortMax.value = row.max
  modalPortToken.value = row.token
  modalPortUser.value = row.user
  createPort.value = false
  modalShow.value = true
}

function deletePort(row) {
  dialog.info({
    title: '是否删除',
    content: '是否删除这条端口规则',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      const res = await axios.delete(`/api/admin/rule/${row.id}`, {
        headers: {
          Authorization: `Bearer ${getToken()}`,
        },
      })
      const data: Response = res.data
      if (data.status === 200) {
        await loadPort()
      } else {
        dialogError(dialog, '删除失败', data.msg)
      }
    },
  })
}

async function modalSubmit() {
  if (createPort.value) {
    const res = await axios.post(
      '/api/admin/rule',
      {
        min: modalPortMin.value,
        max: modalPortMax.value,
        name: modalPortToken.value,
        user: modalPortUser.value,
      },
      {
        headers: {
          Authorization: 'Bearer ' + getToken(),
        },
      },
    )
    const data: Response = res.data
    if (data.status === 200) {
      await loadPort()
    } else {
      dialogError(dialog, '创建失败', data.msg)
    }
  } else {
    const res = await axios.put(
      '/api/admin/rule',
      {
        id: modalPortId.value,
        min: modalPortMin.value,
        max: modalPortMax.value,
        name: modalPortToken.value,
        user: modalPortUser.value,
      },
      {
        headers: {
          Authorization: 'Bearer ' + getToken(),
        },
      },
    )
    const data: Response = res.data
    if (data.status === 200) {
      await loadPort()
    } else {
      dialogError(dialog, '修改失败', data.msg)
    }
  }
  modalShow.value = false
}

async function loadPort() {
  const res = await axios.get('/api/admin/rule', {
    headers: {
      Authorization: 'Bearer ' + getToken(),
    },
  })
  const data: Response = res.data
  if (data.status === 200) {
    ports.ports = data.data.ports
  } else if (data.status == 401) {
    router.push('/login')
  } else if (data.status == 403) {
    router.push('/')
  } else {
    dialogError(dialog, '数据获取失败', data.msg)
  }
}

onMounted(async () => {
  await loadPort()
})
</script>

<template>
  <n-card>
    <template #header>
      <h3>管理员 - 规则管理</h3>
    </template>
    <n-flex>
      <n-button @click="openCreateModal()" size="large" type="primary">创建端口规则</n-button>
      <n-data-table :columns="columns" :data="ports.ports"></n-data-table>
    </n-flex>
  </n-card>
  <n-modal v-model:show="modalShow">
    <n-card style="width: 400px">
      <template #header>
        <h3>{{ createPort ? '创建' : '编辑' }}端口规则</h3>
      </template>
      <n-form>
        <n-form-item label="ID" v-if="!createPort">
          <n-input readonly disabled v-model:value="modalPortId"></n-input>
        </n-form-item>
        <n-form-item label="最小端口">
          <n-input-number v-model:value="modalPortMin"></n-input-number>
        </n-form-item>
        <n-form-item label="最大端口">
          <n-input-number v-model:value="modalPortMax"></n-input-number>
        </n-form-item>
        <n-form-item label="Token">
          <n-input v-model:value="modalPortToken"></n-input>
        </n-form-item>
        <n-form-item label="用户名">
          <n-input v-model:value="modalPortUser"></n-input>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button size="large" @click="modalShow = false">取消</n-button>
          <n-button size="large" type="primary" @click="modalSubmit()">确定</n-button>
        </n-space>
      </template>
    </n-card>
  </n-modal>
</template>

<style scoped></style>
