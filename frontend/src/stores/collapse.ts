import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useCollapseStore = defineStore('collapse', {
  state() {
    const collapse = ref(true)
    return {
      collapse,
    }
  },
  persist: true
})
