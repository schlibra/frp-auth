import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const usePortsStore = defineStore("ports", {
  state() {
    const ports = ref([])
    const count = computed(() => ports.value.length)
    return {
      ports,
      count,
    }
  },
  persist: true
})
