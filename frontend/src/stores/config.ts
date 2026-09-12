import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useConfigStore = defineStore("config", {
  state() {
    const preConfig = ref("")
    const proxyName = ref("")
    const proxyType = ref("")
    const localIP = ref("")
    const localPort = ref("")
    const remotePort = ref("")
    const result = computed(() =>
      proxyName.value && proxyType.value && localIP.value && localPort.value && remotePort.value
        ? `${preConfig.value}

[[proxies]]
name = "${proxyName.value}"
type = "${proxyType.value}"
localIP = "${localIP.value}"
localPort = "${localPort.value}"
remotePort = "${remotePort.value}"
`
        : '',
    )
    const clear = () => {
      preConfig.value = ""
      proxyName.value = ""
      proxyType.value = ""
      localIP.value = ""
      localPort.value = ""
      remotePort.value = ""
    }
    return {
      preConfig,
      proxyName,
      proxyType,
      localIP,
      localPort,
      remotePort,
      result,
      clear,
    }
  },
  persist: true,
})
