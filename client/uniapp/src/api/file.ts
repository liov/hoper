import { httpclient } from '@/api/common'
import { tokenKey } from '@/store/user'
import { PreUploadResp, PreUploadReq, PreUploadType } from '@gen/pb/file/file.service'
import { getMD5 } from '@/utils/file'

type UploadResult = {
  fileUrl: string
}

class FileService {
  static async chooseAndUploadImage(chooseOptions?: UniApp.ChooseImageOptions): Promise<UploadResult> {
    const chooseResult = await uni.chooseImage(chooseOptions || {})
    const tempFilePath = chooseResult.tempFilePaths?.[0] || ''
    if (!tempFilePath) throw new Error('empty tempFilePath')
    const tempFiles = Array.isArray(chooseResult.tempFiles) ? chooseResult.tempFiles : []
    const tempFile = tempFiles[0] as { name?: string; size?: number } | undefined
    return this.uploadImageByPreUpload(tempFilePath, tempFile?.name, tempFile?.size)
  }

  static async uploadImageByPreUpload(tempFilePath: string, name?: string, size?: number): Promise<UploadResult> {
    const fileName = name || tempFilePath.split('/').pop() || `avatar_${Date.now()}.jpg`
    const fileSize = Number(size || 0)
    const data = await this.preUpload({
      preUploadType: PreUploadType.PRE_UPLOAD_TYPE_UPLOAD_URL, name: fileName, size: fileSize,
      md5: await getMD5(tempFilePath),
    })
    const uploadUrl = data?.uploadUrl || ''
    if (!uploadUrl) throw new Error('empty preUpload uploadUrl')
    const uploadResp = await uni.uploadFile({
      url: String(uploadUrl),
      filePath: tempFilePath,
      name: 'file',
      header: { Authorization: uni.getStorageSync(tokenKey) || '' },
    })
    const raw = typeof uploadResp.data === 'string' ? JSON.parse(uploadResp.data || '{}') : uploadResp.data || {}
    const fileUrl = data?.file?.url || raw?.data?.url || raw?.data?.path || raw?.url || raw?.path || tempFilePath
    if (!fileUrl) throw new Error('empty upload response')
    return { fileUrl: String(fileUrl) }
  }
  static async preUpload(req: PreUploadReq): Promise<PreUploadResp> {
    return await httpclient.post<PreUploadResp>('/api/preUpload', { data: req, decode: PreUploadResp })
  }
}

export default FileService
