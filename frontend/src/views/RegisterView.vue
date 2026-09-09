<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import router from '@/router'
import HeaderComponent from '@/components/header-component.vue'
import type Response from '@/model/response.ts'
import { dialogError, dialogSuccess } from '@/utils/dialog.ts'
import { useDialog } from 'naive-ui'

const dialog = useDialog()
const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const nickname = ref('')
const passwordRef = ref(null)
const confirmPasswordRef = ref(null)
const nicknameRef = ref(null)

const doRegister = async () => {
  if (password.value == confirmPassword.value) {
    const result = await axios.post('/api/user/register', {
      username: username.value,
      password: password.value,
      nickname: nickname.value,
    })
    const data: Response = result.data
    if (data.status == 200) {
      dialogSuccess(dialog, '注册成功', data.msg, () => {
        router.push('/login')
      })
    } else {
      dialogError(dialog, '注册失败', data.msg)
    }
  } else {
    dialogError(dialog, "注册失败", "两次密码不一致", ()=>{
      confirmPasswordRef.value.focus()
    })
  }
}
</script>

<template>
  <n-layout>
    <n-layout-header>
      <header-component title="注册"></header-component>
    </n-layout-header>
    <n-layout-content>
      <n-flex justify="center" align="center" style="height: calc(100vh - 200px)">
        <n-card size="large" style="max-width: 600px">
          <template #header>
            <span>注册账号</span>
          </template>
          <n-form>
            <n-form-item label="用户名">
              <n-input v-model:value="username" @keydown.enter="passwordRef.focus()"></n-input>
            </n-form-item>
            <n-form-item label="密码">
              <n-input
                ref="passwordRef"
                type="password"
                v-model:value="password"
                @keydown.enter="confirmPasswordRef.focus()"
              ></n-input>
            </n-form-item>
            <n-form-item label="确认密码">
              <n-input
                ref="confirmPasswordRef"
                type="password"
                v-model:value="confirmPassword"
                @keydown.enter="nicknameRef.focus()"
              ></n-input>
            </n-form-item>
            <n-form-item label="昵称">
              <n-input
                ref="nicknameRef"
                v-model:value="nickname"
                @keydown.enter="doRegister()"
              ></n-input>
            </n-form-item>
          </n-form>
          <template #footer>
            <n-button @click="doRegister()" size="large" type="primary">注册</n-button>
          </template>
        </n-card>
      </n-flex>
    </n-layout-content>
    <n-layout-footer></n-layout-footer>
  </n-layout>
</template>

<style scoped></style>
