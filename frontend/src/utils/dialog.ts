import type { DialogApiInjection } from 'naive-ui/es/dialog/src/DialogProvider'

function dialogSuccess(dialog: DialogApiInjection, title: string, content: string, callback: () => void = () => {}) {
  dialog.success({
    title,
    content,
    positiveText: '确定',
    onPositiveClick() {
      callback()
    },
    onClose() {
      callback()
    },
    onEsc() {
      callback()
    },
  })
}
function dialogError(dialog: DialogApiInjection, title: string, content: string, callback: () => void = () => {}) {
  dialog.error({
    title,
    content,
    closable: false,
    closeFocusable: false,
    closeOnEsc: false,
    positiveText: '确定',
    onPositiveClick() {
      callback()
    },
  })
}

export { dialogSuccess, dialogError }
