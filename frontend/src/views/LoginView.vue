<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import router from '@/router'
import HeaderComponent from '@/components/header-component.vue'
import { setToken } from '@/utils/token.ts'
import type LoginToken from '@/model/loginToken.ts'
import type Response from '@/model/response.ts'
import { dialogError, dialogSuccess } from '@/utils/dialog.ts'
import { useDialog, useMessage } from 'naive-ui'
import JSEncrypt from 'jsencrypt'

const dialog = useDialog()
const message = useMessage()
const username = ref('')
const password = ref('')
const passwordRef = ref(null)
const loading = ref(false)

function encrypt(data: string, pubKey: string) {
  const e = new JSEncrypt()
  e.setPublicKey(pubKey)

  const ed = e.encrypt(data)
  if (!ed) {
    dialogError(dialog, '加密失败', '数据加密失败')
  }
  return ed
}

async function doLogin() {
  loading.value = true
  const res = await axios.put('/api/user/login', {
    username: username.value,
  })
  let data: Response = res.data
  if (data.status == 200) {
    console.log(data)
  } else {
    dialogError(dialog, '登录失败', data.msg)
    loading.value = false
    return
  }
  const _p = encrypt(password.value, data.data.public_key)
  const result = await axios.post('/api/user/login', {
    username: username.value,
    password: _p,
  })
  data = result.data
  if (data.status == 200) {
    message.success('登录成功', {
      closable: true,
      duration: 5000,
    })
    setToken(data.data.token)
    setTimeout(() => {
      router.push('/')
    }, 500)
    // dialogSuccess(dialog, '登录成功', data.msg, () => {
    //   const loginData: LoginToken = data.data
    //   setToken(loginData.token)
    //   router.push('/')
    // })
  } else {
    dialogError(dialog, '登录失败', data.msg)
  }
  loading.value = false
}
</script>

<template>
  <n-layout>
    <n-layout-header>
      <header-component title="登录"></header-component>
    </n-layout-header>
    <n-layout-content>
      <n-flex justify="center" align="center" style="height: calc(100vh - 200px)">
        <n-card size="large" style="max-width: 600px">
          <template #header>
            <span>登录账号</span>
          </template>
          <n-form>
            <n-form-item label="用户名">
              <n-input v-model:value="username" @keydown.enter="passwordRef.focus()"></n-input>
            </n-form-item>
            <n-form-item label="密码" @keydown.enter="doLogin()">
              <n-input ref="passwordRef" type="password" v-model:value="password"></n-input>
            </n-form-item>
          </n-form>
          <template #footer>
            <n-button @click="doLogin()" size="large" type="primary" :loading="loading">登录</n-button>
          </template>
        </n-card>
      </n-flex>
    </n-layout-content>
    <n-layout-footer></n-layout-footer>
  </n-layout>
</template>

<style scoped></style>
