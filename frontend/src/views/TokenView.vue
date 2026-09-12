<script setup lang="ts">
import HeaderComponent from '@/components/header-component.vue'
import MenuComponent from '@/components/menu-component.vue'
import { h, onMounted, ref } from 'vue'
import axios from 'axios'
import type Response from '@/model/response.ts'
import { dialogError } from '@/utils/dialog.ts'
import { type DataTableColumns, NButton, NSpace, NSwitch, useDialog } from 'naive-ui'
import { getToken } from '@/utils/token.ts'
import router from '@/router'

const dialog = useDialog()
const tokenList = ref([])
const tokenModalShow = ref(false)
const newTokenName = ref('')
const frpConfig = ref('')
const configModalShow = ref(false)
const frpConfigText = ref(null)

const columns: DataTableColumns = [
  {
    title: 'ID',
    key: 'id',
    width: 50,
  },
  {
    title: '名称',
    key: 'name',
    width: 100,
  },
  {
    title: 'Token',
    key: 'token',
    width: 200,
  },
  {
    title: '用户',
    key: 'user',
    width: 100,
  },
  {
    title: '是否启用',
    key: 'enable',
    render(row) {
      return h(NSwitch, {
        value: row.enable == 1,
        onUpdateValue(value: boolean) {
          changeTokenEnable(row.name, value)
        },
      })
    },
    width: 100,
  },
  {
    title: '操作',
    key: 'action',
    minWidth: 400,
    render(row) {
      return h(NSpace, {}, [
        h(
          NButton,
          {
            type: 'primary',
            onClick() {
              generateConf(row.id)
            },
          },
          '生成配置',
        ),
        h(
          NButton,
          {
            type: 'warning',
            onClick() {
              regenerateToken(row.name)
            },
          },
          { default: () => '重新生成Token' },
        ),
        h(
          NButton,
          {
            type: 'error',
            onClick() {
              deleteToken(row.name)
            },
          },
          { default: () => '删除Token' },
        ),
      ])
    },
  },
]

async function regenerateToken(name: string) {
  dialog.info({
    title: '是否重新生成Token',
    content: '是否重新生成Token，原Token将失效',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      const res = await axios.patch(
        '/api/token',
        {
          name,
        },
        {
          headers: {
            Authorization: 'Bearer ' + getToken(),
          },
        },
      )
      const data: Response = res.data
      if (data.status == 200) {
        await loadTokens()
      } else {
        dialogError(dialog, 'Token修改失败', data.msg)
      }
    },
  })
}

async function deleteToken(name: string) {
  dialog.info({
    title: '是否删除',
    content: '是否删除Token？',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      const res = await axios.delete(`/api/token?name=${name}`, {
        headers: {
          Authorization: 'Bearer ' + getToken(),
        },
      })
      const data: Response = res.data
      if (data.status == 200) {
        await loadTokens()
      } else {
        dialogError(dialog, '删除失败', data.msg)
      }
    },
  })
}
async function changeTokenEnable(name: string, status: boolean) {
  const res = await axios.put(
    '/api/token',
    {
      name: name,
      enable: status,
    },
    {
      headers: {
        Authorization: 'Bearer ' + getToken(),
      },
    },
  )
  const data: Response = res.data
  if (data.status == 200) {
    await loadTokens()
  } else {
    dialogError(dialog, '状态修改失败', data.msg)
  }
}

async function generateConf(id: number) {
  const res = await axios.get(`/api/token/conf/${id}`, {
    headers: {
      Authorization: 'Bearer ' + getToken(),
    },
  })
  const data: Response = res.data
  if (data.status == 200) {
    frpConfig.value = data.data.config
    configModalShow.value = true
  } else {
    dialogError(dialog, '配置生成失败', data.msg)
  }
}

async function loadTokens() {
  const res = await axios.get('/api/token', {
    headers: {
      Authorization: 'Bearer ' + getToken(),
    },
  })
  const data: Response = res.data
  if (data.status == 200) {
    tokenList.value = data.data.tokens
  } else if (data.status == 401) {
    router.push('/login')
  } else {
    dialogError(dialog, '数据获取失败', data.msg)
  }
}
async function createToken() {
  const res = await axios.post(
    '/api/token',
    {
      name: newTokenName.value,
    },
    {
      headers: {
        Authorization: 'Bearer ' + getToken(),
      },
    },
  )
  const data: Response = res.data
  if (data.status == 200) {
    await loadTokens()
  } else {
    dialogError(dialog, '创建失败', data.msg)
  }
  newTokenName.value = ''
  tokenModalShow.value = false
}

function copyConfig() {
  frpConfigText.value.select()
  document.execCommand("copy")
}

onMounted(async () => {
  await loadTokens()
})
</script>

<template>
  <n-card>
    <template #header>
      <h3>Token管理</h3>
    </template>
    <n-space vertical>
      <n-button size="large" type="primary" @click="tokenModalShow = true"
        >创建Token</n-button
      >
      <n-data-table :columns="columns" :data="tokenList"> </n-data-table>
    </n-space>
  </n-card>
  <n-modal v-model:show="tokenModalShow">
    <n-card style="width: 400px">
      <template #header>
        <h3>创建Token</h3>
      </template>
      <n-form>
        <n-form-item label="Token名称">
          <n-input v-model:value="newTokenName"></n-input>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button size="large" @click="tokenModalShow = false">取消</n-button>
          <n-button size="large" type="primary" @click="createToken()">确定</n-button>
        </n-space>
      </template>
    </n-card>
  </n-modal>
  <n-modal v-model:show="configModalShow">
    <n-card style="width: 600px">
      <template #header>
        <h3>Frp配置</h3>
      </template>
      <n-form>
        <n-form-item label="Frp配置">
          <n-input ref="frpConfigText" :rows="5" :value="frpConfig" type="textarea"></n-input>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button type="info" size="large" @click="copyConfig()">复制</n-button>
          <n-button type="primary" size="large" @click="configModalShow = false">确定</n-button>
        </n-space>
      </template>
    </n-card>
  </n-modal>
</template>

<style scoped></style>
