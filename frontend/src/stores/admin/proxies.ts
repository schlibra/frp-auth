import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useAdminProxiesStore = defineStore("adminProxies", {
  state() {
    const proxies = ref([])
    const count = computed(() => proxies.value.length)
    return {
      proxies,
      count,
    }
  },
  persist: true
})
