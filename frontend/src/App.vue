<script setup lang="ts">
import { darkTheme, zhCN } from 'naive-ui'
import HeaderComponent from '@/components/header-component.vue'
import MenuComponent from '@/components/menu-component.vue'
import { computed, ref } from 'vue'

const path = ref(location.pathname)
const showMenu = computed(() => path.value != '/login' && path.value != '/register' && path.value != '/404')
setInterval(() => {
  path.value = location.pathname
}, 100)
</script>

<template>
  <n-config-provider :theme="darkTheme" :locale="zhCN">
    <n-dialog-provider>
      <n-modal-provider>
        <n-message-provider>
          <n-layout class="main">
            <n-layout-header>
              <header-component></header-component>
            </n-layout-header>
            <n-layout has-sider>
              <menu-component v-if="showMenu"></menu-component>
              <n-layout-content>
                <n-flex justify="center" align="center">
                  <router-view></router-view>
                </n-flex>
              </n-layout-content>
            </n-layout>
          </n-layout>
        </n-message-provider>
      </n-modal-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>

<style scoped>
.main {
  height: 100vh;
}
</style>
