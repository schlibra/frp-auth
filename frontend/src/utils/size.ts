function calcSize(size: number) {
  const suffix = ["B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"]
  let idx = 0
  while (size > 1024) {
    idx+=1
    size = Math.round(size / 1024 * 10) / 10
  }
  return `${size}${suffix[idx]}`
}

export {
  calcSize
}
