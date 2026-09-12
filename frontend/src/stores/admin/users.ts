import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useAdminUsersStore = defineStore("adminUsers", {
  state() {
    const users = ref([])
    const count = computed(() => users.value.length)
    return {
      users,
      count
    }
  },
  persist: true
})
