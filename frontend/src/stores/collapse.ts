import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useCollapseStore = defineStore('collapse', {
  state() {
    const collapse = ref(false)
    return {
      collapse,
    }
  },
  persist: true
})
