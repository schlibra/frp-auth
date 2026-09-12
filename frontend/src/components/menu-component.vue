<script setup lang="ts">
import { type MenuOption, NIcon, useDialog } from 'naive-ui'
import { type Component, computed, type ComputedRef, h, onMounted, ref } from 'vue'
import {
  HomeOutline as HomeIcon,
  KeyOutline as KeyIcon,
  GlobeOutline as PortIcon,
  PersonOutline as UserIcon,
  ServerOutline as AdminIcon,
  DesktopOutline as ClientIcon,
  RepeatOutline as ProxyIcon,
  SaveOutline as ConfigIcon,
} from '@vicons/ionicons5'
import axios from 'axios'
import router from '@/router'
import type Response from '@/model/response'
import type UserInfo from '@/model/userInfo.ts'
import { useCollapseStore } from '@/stores/collapse.ts'

const renderIcon = (icon: Component) => () => h(NIcon, null, { default: () => h(icon) })
const admin = ref(false)
const dialog = useDialog()
const collapsed = useCollapseStore()
const urlPath = ref(location.pathname)

const menuOptions: ComputedRef<MenuOption[]> = computed(() => [
  {
    label: '首页',
    key: '/',
    icon: renderIcon(HomeIcon),
  },
  {
    label: 'Token管理',
    key: '/token',
    icon: renderIcon(KeyIcon),
  },
  {
    label: '端口规则管理',
    key: '/port',
    icon: renderIcon(PortIcon),
  },
  {
    label: '客户端管理',
    key: '/client',
    icon: renderIcon(ClientIcon),
  },
  {
    label: '映射管理',
    key: '/proxy',
    icon: renderIcon(ProxyIcon),
  },
  {
    label: "配置生成",
    key: '/config',
    icon: renderIcon(ConfigIcon),
  },
  {
    label: '个人设置',
    key: '/user',
    icon: renderIcon(UserIcon),
  },
  {
    label: '管理员设置',
    icon: renderIcon(AdminIcon),
    show: admin.value,
    children: [
      {
        label: '用户管理',
        key: '/admin/user',
        icon: renderIcon(UserIcon),
      },
      {
        label: 'Token管理',
        key: '/admin/token',
        icon: renderIcon(KeyIcon),
      },
      {
        label: '端口规则管理',
        key: '/admin/port',
        icon: renderIcon(PortIcon),
      },
      {
        label: '客户端管理',
        key: '/admin/client',
        icon: renderIcon(ClientIcon),
      },
      {
        label: '映射管理',
        key: '/admin/proxy',
        icon: renderIcon(ProxyIcon),
      }
    ],
  },
])


onMounted(async () => {

  const token = localStorage.getItem('token')
  const res = await axios.get('/api/user/', {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })
  const data: Response = res.data
  if (data.status == 200) {
    const userData: UserInfo = data.data
    admin.value = userData.role === 'admin'
  } else if (data.status == 401) {
    router.push('/login')
  } else {
    dialog.error({
      title: '数据获取失败',
      content: data.msg,
      positiveText: '确定',
      closable: false,
      closeOnEsc: false,
      closeFocusable: false,
      onPositiveClick() {
        location.reload()
      },
    })
  }
})

const menuChange = (key: string) => {
  router.push(key)
}
setInterval(() => {
  urlPath.value = location.pathname
}, 100)
</script>

<template>
  <n-layout-sider
    bordered
    collapse-mode="width"
    :collapsed-width="48"
    :width="240"
    :collapsed="collapsed.collapse"
    @collapse="collapsed.collapse = true"
    @expand="collapsed.collapse = false"
    show-trigger
  >
    <n-menu
      :value="urlPath"
      class="menu"
      :options="menuOptions"
      :collapsed="collapsed.collapse"
      :collapsed-width="48"
      :collapsed-icon-siz="20"
      @update-value="menuChange"
    ></n-menu>
  </n-layout-sider>
</template>

<style scoped>
.menu {
  height: calc(100vh - 75px);
}
</style>
