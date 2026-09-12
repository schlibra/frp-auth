import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useProxiesStore = defineStore("proxies", {
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
