import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useAdminPortsStore = defineStore("adminPorts", {
  state() {
    const ports = ref([])
    const count = computed(() => ports.value.length)
    return {
      ports,
      count
    }
  },
  persist: true
})
