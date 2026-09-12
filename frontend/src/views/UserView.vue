<script setup lang="ts">
import HeaderComponent from '@/components/header-component.vue'
import MenuComponent from '@/components/menu-component.vue'
import { onMounted, ref } from 'vue'
import axios from 'axios'
import type Response from '@/model/response.ts'
import { dialogError, dialogSuccess } from '@/utils/dialog.ts'
import router from '@/router'
import type UserInfo from '@/model/userInfo.ts'
import { useDialog } from 'naive-ui'
import { removeToken } from '@/utils/token.ts'

const username = ref('')
const nickname = ref('')
const password = ref('')
const role = ref('')
const dialog = useDialog()

async function loadUserInfo() {
  const res = await axios.get('/api/user', {
    headers: {
      Authorization: 'Bearer ' + localStorage.getItem('token'),
    },
  })
  const data: Response = res.data
  if (data.status == 200) {
    const userData: UserInfo = data.data
    username.value = userData.username
    nickname.value = userData.nickname
    role.value = userData.role == 'admin' ? '管理员' : '普通用户'
  } else if (data.status == 401) {
    router.push('/login')
  } else {
    dialogError(dialog, '数据获取失败', data.msg, () => {
      router.push('/login')
    })
  }
}

async function updateUserInfo() {
  const res = await axios.put(
    '/api/user',
    {
      nickname: nickname.value,
    },
    {
      headers: {
        Authorization: 'Bearer ' + localStorage.getItem('token'),
      },
    },
  )
  const data: Response = res.data
  if (data.status == 200) {
    dialogSuccess(dialog, '修改成功', data.msg, async () => {
      await loadUserInfo()
    })
  } else {
    dialogError(dialog, '修改失败', data.msg)
  }
}

async function changeUserPassword() {
  if (password.value === '') {
    dialogError(dialog, '修改失败', '密码不能为空')
  } else {
    const res = await axios.patch(
      '/api/user',
      {
        password: password.value,
      },
      {
        headers: {
          Authorization: 'Bearer ' + localStorage.getItem('token'),
        },
      },
    )
    const data: Response = res.data
    if (data.status == 200) {
      dialogSuccess(dialog, '密码修改成功', '密码修改成功，请重新登录', () => {
        removeToken()
        router.push('/login')
      })
    }
  }
}

async function logoutUser() {
  dialog.info({
    title: '退出登录',
    content: '是否退出登录？',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      await axios.post(
        '/api/user/logout',
        {},
        {
          headers: {
            Authorization: 'Bearer ' + localStorage.getItem('token'),
          },
        },
      )
      removeToken()
      router.push('/login')
    },
  })
}

onMounted(async () => {
  await loadUserInfo()
})
</script>

<template>
  <n-card>
    <template #header>
      <h3>个人设置</h3>
    </template>
    <n-form>
      <n-form-item label="用户名">
        <n-input readonly disabled placeholder="用户名" v-model:value="username"></n-input>
      </n-form-item>
      <n-form-item label="昵称">
        <n-input placeholder="昵称" v-model:value="nickname"></n-input>
        <n-button @click="updateUserInfo()">修改昵称</n-button>
      </n-form-item>
      <n-form-item label="密码">
        <n-input placeholder="新密码" type="password" v-model:value="password"></n-input>
        <n-button @click="changeUserPassword()">修改密码</n-button>
      </n-form-item>
      <n-form-item label="身份">
        <n-input readonly disabled v-model:value="role"></n-input>
      </n-form-item>
    </n-form>
    <template #footer>
      <n-button size="large" type="error" @click="logoutUser()">退出登录</n-button>
    </template>
  </n-card>
</template>

<style scoped></style>
