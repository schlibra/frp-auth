import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useTokensStore = defineStore('tokens', {
  state() {
    const tokens = ref([])
    const count = computed(()=>tokens.value.length)
    return {
      tokens,
      count
    }
  },
  persist: true
})
