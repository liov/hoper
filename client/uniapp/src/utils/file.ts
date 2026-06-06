import SparkMD5 from 'spark-md5'

const SparkMD5Any = SparkMD5 as any

export async function getMD5(tempFilePath: string): Promise<string> {
  const platform = uni.getSystemInfoSync().uniPlatform
  if (platform === 'web') {
    const response = await fetch(tempFilePath)
    const arrayBuffer = await response.arrayBuffer()
    return SparkMD5Any.ArrayBuffer.hash(arrayBuffer)
  }
  const fileSystemManager = (uni as any).getFileSystemManager?.()
  if (!fileSystemManager) {
    throw new Error('current platform does not support fileSystemManager')
  }
  const readResult = await new Promise<UniApp.ReadFileSuccessCallbackResult>((resolve, reject) => {
    fileSystemManager.readFile({
      filePath: tempFilePath,
      success: resolve,
      fail: reject,
    })
  })
  const arrayBuffer = readResult.data as ArrayBuffer
  return SparkMD5Any.ArrayBuffer.hash(arrayBuffer)
}
