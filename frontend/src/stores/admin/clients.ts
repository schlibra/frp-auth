import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useAdminClientsStore = defineStore("adminClients", {
  state() {
    const clients = ref([])
    const count = computed(() => clients.value.length)
    return {
      clients,
      count,
    }
  },
  persist: true,
})
