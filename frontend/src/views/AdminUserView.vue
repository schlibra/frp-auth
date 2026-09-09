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

const userList = ref([])
const dialog = useDialog()
const modalShow = ref(false)
const createUser = ref(false)
const modalUserId = ref(0)
const modalUsername = ref('')
const modalPassword = ref('')
const modalNickname = ref('')
const modalAdmin = ref(false)
const modalEnable = ref(false)

const columns: Ref<DataTableColumn[]> = ref([
  {
    title: 'ID',
    key: 'id',
    width: 50
  },
  {
    title: '用户名',
    key: 'username',
    width: 120
  },
  {
    title: '昵称',
    key: 'nickname',
    width: 120
  },
  {
    title: '管理员',
    key: 'admin',
    width: 70,
    render(row) {
      return h(NSwitch, {
        value: row.admin == 1,
        async onUpdateValue(status: boolean) {
          await changeAdmin(row, status)
        },
      })
    },
  },
  {
    title: '已启用',
    key: 'enable',
    width: 70,
    render: (row) => {
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
    width: 150,
    render: (row) => {
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
              deleteUser(row)
            },
          },
          '删除',
        ),
      ])
    },
  },
])

function openCreateModal() {
  modalUserId.value = 0
  modalUsername.value = ''
  modalNickname.value = ''
  modalPassword.value = ''
  modalAdmin.value = false
  modalEnable.value = true
  createUser.value = true
  modalShow.value = true
}

function openUpdateModal(row) {
  modalUserId.value = row.id
  modalUsername.value = row.username
  modalNickname.value = row.nickname
  modalPassword.value = ''
  modalAdmin.value = row.admin
  modalEnable.value = row.enable
  createUser.value = false
  modalShow.value = true
}

async function loadUser() {
  const res = await axios.get('/api/admin/user', {
    headers: {
      Authorization: 'Bearer ' + getToken(),
    },
  })
  const data: Response = res.data
  if (data.status === 200) {
    userList.value = data.data.users
  } else if (data.status == 401) {
    router.push('/login')
  } else if (data.status == 403) {
    router.push('/')
  } else {
    dialogError(dialog, '数据获取失败', data.msg)
  }
}

async function modalSubmit() {
  if (createUser.value) {
    const res = await axios.post(
      '/api/admin/user',
      {
        username: modalUsername.value,
        password: modalPassword.value,
        nickname: modalNickname.value,
        admin: modalAdmin.value,
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
      await loadUser()
    } else {
      dialogError(dialog, '创建失败', data.msg)
    }
  } else {
    const res = await axios.put(
      '/api/admin/user',
      {
        username: modalUsername.value,
        password: modalPassword.value,
        nickname: modalNickname.value,
        admin: modalAdmin.value,
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
      await loadUser()
    } else {
      dialogError(dialog, '修改失败', data.msg)
    }
  }
  modalShow.value = false
}

async function changeAdmin(row, status: boolean) {
  row.admin = status
  const res = await axios.put('/api/admin/user', row, {
    headers: {
      Authorization: 'Bearer ' + getToken(),
    },
  })
  const data: Response = res.data
  if (data.status === 200) {
    await loadUser()
  } else {
    dialogError(dialog, '修改失败', data.msg)
    await loadUser()
  }
}

async function changeEnable(row, status: boolean) {
  row.enable = status
  const res = await axios.put('/api/admin/user', row, {
    headers: {
      Authorization: 'Bearer ' + getToken(),
    },
  })
  const data: Response = res.data
  if (data.status === 200) {
    await loadUser()
  } else {
    dialogError(dialog, '修改失败', data.msg)
    await loadUser()
  }
}

async function deleteUser(row) {
  dialog.info({
    title: '是否删除',
    content: '是否删除该用户',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      const res = await axios.delete(`/api/admin/user/${row.id}`, {
        headers: {
          Authorization: 'Bearer ' + getToken(),
        },
      })
      const data: Response = res.data
      if (data.status === 200) {
        await loadUser()
      } else {
        dialogError(dialog, '删除失败', data.msg)
      }
    },
  })
}

onMounted(async () => {
  await loadUser()
})
</script>

<template>
  <n-layout>
    <n-layout-header>
      <header-component title="管理员 - 用户管理"></header-component>
    </n-layout-header>
    <n-layout has-sider>

        <menu-component></menu-component>

      <n-layout-content>
        <n-card>
          <template #header>
            <h3>管理员 - 用户管理</h3>
          </template>
          <n-space vertical>
            <n-button @click="openCreateModal()" size="large" type="primary">创建用户</n-button>
            <n-data-table :columns="columns" :data="userList"> </n-data-table>
          </n-space>
        </n-card>
      </n-layout-content>
    </n-layout>
  </n-layout>
  <n-modal v-model:show="modalShow">
    <n-card style="width: 400px">
      <template #header>
        <h3>{{ createUser ? '创建' : '编辑' }}用户</h3>
      </template>
      <n-form>
        <n-form-item label="ID" v-if="!createUser">
          <n-input readonly disabled v-model:value="modalUserId"></n-input>
        </n-form-item>
        <n-form-item label="用户名">
          <n-input
            :readonly="!createUser"
            :disabled="!createUser"
            v-model:value="modalUsername"
          ></n-input>
        </n-form-item>
        <n-form-item label="密码">
          <n-input type="password" v-model:value="modalPassword"></n-input>
        </n-form-item>
        <n-form-item label="昵称">
          <n-input v-model:value="modalNickname"></n-input>
        </n-form-item>
        <n-form-item label="管理员" v-if="createUser">
          <n-switch v-model:value="modalAdmin"></n-switch>
        </n-form-item>
        <n-form-item label="已启用" v-if="createUser">
          <n-switch v-model:value="modalEnable"></n-switch>
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
