<script setup lang="ts">
import HeaderComponent from '@/components/header-component.vue'
import MenuComponent from '@/components/menu-component.vue'
import { h, onMounted, ref } from 'vue'
import axios from 'axios'
import { getToken } from '@/utils/token.ts'
import type Response from '@/model/response.ts'
import { type DataTableColumn, NButton, NSpace, useDialog } from 'naive-ui'
import { dialogError } from '@/utils/dialog.ts'
import router from '@/router'

const dialog = useDialog()
const portData = ref([])
const modalShow = ref(false)
const createPort = ref(false)
const modalPortId = ref(0)
const modalPortMin = ref(0)
const modalPortMax = ref(0)
const modalPortToken = ref('')

const columns: DataTableColumn[] = [
  {
    title: 'ID',
    key: 'id',
    width: 50,
  },
  {
    title: '最小端口',
    key: 'min',
    width: 80,
  },
  {
    title: '最大端口',
    key: 'max',
    width: 80
  },
  {
    title: 'Token',
    key: 'token',
    width: 120,
  },
  {
    title: '用户',
    key: 'user',
    width: 120
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
              updatePort(row)
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
]

function updatePort(row) {
  modalPortId.value = row.id
  modalPortMin.value = row.min
  modalPortMax.value = row.max
  modalPortToken.value = row.token
  modalShow.value = true
  createPort.value = false
}

function deletePort(row) {
  dialog.info({
    title: '是否删除',
    content: '是否删除这条端口规则',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      const res = await axios.delete(`/api/rule/${row.id}`, {
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

function openCreateModal() {
  modalPortId.value = 0
  modalPortMin.value = ""
  modalPortMax.value = ""
  modalPortToken.value = ''
  createPort.value = true
  modalShow.value = true
}

async function modalSubmit() {
  if (createPort.value) {
    const res = await axios.post(
      '/api/rule',
      {
        name: modalPortToken.value,
        min: modalPortMin.value,
        max: modalPortMax.value,
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
      '/api/rule',
      {
        id: modalPortId.value,
        name: modalPortToken.value,
        min: modalPortMin.value,
        max: modalPortMax.value,
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
  const res = await axios.get('/api/rule', {
    headers: {
      Authorization: 'Bearer ' + getToken(),
    },
  })
  const data: Response = res.data
  if (data.status == 200) {
    portData.value = data.data.ports
  } else if (data.status == 401) {
    router.push('/login')
  } else {
    dialogError(dialog, '数据获取失败', data.msg)
  }
}

onMounted(async () => {
  await loadPort()
})
</script>

<template>
  <n-layout>
    <n-layout-header>
      <header-component title="规则管理"></header-component>
    </n-layout-header>
    <n-layout has-sider>
      <menu-component></menu-component>
      <n-layout-content>
        <n-card>
          <template #header>
            <h3>规则管理</h3>
          </template>
          <n-flex>
            <n-button @click="openCreateModal()" size="large" type="primary">创建端口规则</n-button>
            <n-data-table :columns="columns" :data="portData"></n-data-table>
          </n-flex>
        </n-card>
      </n-layout-content>
    </n-layout>
  </n-layout>
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
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="modalShow = false">取消</n-button>
          <n-button @click="modalSubmit()">确定</n-button>
        </n-space>
      </template>
    </n-card>
  </n-modal>
</template>

<style scoped></style>
