<script setup lang="ts">
import HeaderComponent from '@/components/header-component.vue'
import MenuComponent from '@/components/menu-component.vue'
import { h, onMounted, type Ref, ref } from 'vue'
import axios from 'axios'
import { getToken } from '@/utils/token.ts'
import type Response from '@/model/response.ts'
import { type DataTableColumn, NButton, NSpace, NSwitch, useDialog } from 'naive-ui'
import { dialogError } from '@/utils/dialog.ts'
import router from '@/router'

const tokenList = ref([])
const dialog = useDialog()
const modalShow = ref(false)
const createToken = ref(false)
const modalName = ref('')
const modalUser = ref('')
const modalEnable = ref(false)
const columns: Ref<DataTableColumn[]> = ref([
  {
    title: 'ID',
    key: 'id',
    width: 50
  },
  {
    title: '名称',
    key: 'name',
    width: 120
  },
  {
    title: 'Token',
    key: 'token',
    width: 170
  },
  {
    title: '用户',
    key: 'user',
    width: 100
  },
  {
    title: '已启用',
    key: 'enable',
    width: 70,
    render(row) {
      return h(NSwitch, {
        value: row.enable == 1,
        async onUpdateValue(status: boolean) {
          await changeEnable(row, status)
        },
      })
    },
  },
  {
    title: '操作',
    key: 'action',
    width: 290,
    render(row) {
      return h(NSpace, {}, [
        h(
          NButton,
          {
            type: 'primary',
            onClick() {
              openUpdateModal(row)
            },
          },
          '编辑',
        ),
        h(
          NButton,
          {
            type: 'warning',
            onClick() {
              regenerateToken(row)
            },
          },
          '重新生成Token',
        ),
        h(
          NButton,
          {
            type: 'error',
            onClick() {
              deleteToken(row)
            },
          },
          '删除',
        ),
      ])
    },
  },
])

function openCreateModal() {
  modalUser.value = ''
  modalName.value = ''
  createToken.value = true
  modalShow.value = true
}
function openUpdateModal(row) {
  modalName.value = row.name
  modalUser.value = row.user
  modalEnable.value = row.enable == 1
  createToken.value = false
  modalShow.value = true
}

async function modalSubmit() {
  if (createToken.value) {
    const res = await axios.post(
      '/api/admin/token',
      {
        name: modalName.value,
        username: modalUser.value,
      },
      {
        headers: {
          Authorization: 'Bearer ' + getToken(),
        },
      },
    )
    const data: Response = res.data
    if (data.status === 200) {
      await loadToken()
    } else {
      dialogError(dialog, '创建失败', data.msg)
    }
  } else {
    const res = await axios.put(
      '/api/admin/token',
      {
        name: modalName.value,
        username: modalUser.value,
        enable: modalEnable.value,
      },
      {
        headers: {
          Authorization: 'Bearer ' + getToken(),
        },
      },
    )
    const data: Response = res.data
    if (data.status === 200) {
      await loadToken()
    } else {
      dialogError(dialog, '修改失败', data.msg)
    }
  }
  modalShow.value = false
}

async function regenerateToken(row) {
  dialog.info({
    title: '是否重新生成Token',
    content: '是否重新生成Token，原Token将失效',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      const res = await axios.patch(
        '/api/admin/token',
        {
          name: row.name,
        },
        {
          headers: {
            Authorization: 'Bearer ' + getToken(),
          },
        },
      )
      const data: Response = res.data
      if (data.status === 200) {
        await loadToken()
      } else {
        dialogError(dialog, '生成失败', data.msg)
      }
    },
  })
}

async function changeEnable(row, status: boolean) {
  row.enable = status
  row.username = row.user
  const res = await axios.put('/api/admin/token', row, {
    headers: {
      Authorization: 'Bearer ' + getToken(),
    },
  })
  const data: Response = res.data
  if (data.status === 200) {
    await loadToken()
  } else {
    dialogError(dialog, '修改失败', data.msg)
  }
}

async function deleteToken(row) {
  dialog.info({
    title: '是否删除',
    content: '是否删除这个Token',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      const res = await axios.delete(`/api/admin/token?name=${row.name}`, {
        headers: {
          Authorization: 'Bearer ' + getToken(),
        },
      })
      const data: Response = res.data
      if (data.status === 200) {
        await loadToken()
      } else {
        dialogError(dialog, '删除失败', data.msg)
      }
    },
  })
}

async function loadToken() {
  const res = await axios.get('/api/admin/token', {
    headers: {
      Authorization: 'Bearer ' + getToken(),
    },
  })
  const data: Response = res.data
  if (data.status === 200) {
    tokenList.value = data.data.tokens
  } else if (data.status == 401) {
    router.push('/login')
  } else if (data.status == 403) {
    router.push('/')
  } else {
    dialogError(dialog, '数据获取失败', data.msg)
  }
}

onMounted(async () => {
  await loadToken()
})
</script>

<template>
  <n-layout>
    <n-layout-header>
      <header-component title="管理员 - Token管理"></header-component>
    </n-layout-header>
    <n-layout has-sider>

        <menu-component></menu-component>

      <n-layout-content>
        <n-card>
          <template #header>
            <h3>管理员 - Token管理</h3>
          </template>
          <n-space vertical>
            <n-button type="primary" size="large" @click="openCreateModal()">创建Token</n-button>
            <n-data-table :columns="columns" :data="tokenList"></n-data-table>
          </n-space>
        </n-card>
      </n-layout-content>
    </n-layout>
  </n-layout>
  <n-modal v-model:show="modalShow">
    <n-card style="width: 400px">
      <template #header>
        <h3>{{ createToken ? '创建' : '编辑' }}Token</h3>
      </template>
      <n-form>
        <n-form-item label="名称">
          <n-input
            :readonly="!createToken"
            :disabled="!createToken"
            v-model:value="modalName"
          ></n-input>
        </n-form-item>
        <n-form-item label="用户名">
          <n-input v-model:value="modalUser"></n-input>
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
