// 全局要用的类型放到这里

// uni.uploadFile文件上传参数
type IUniUploadFileOptions = {
  file?: File
  files?: UniApp.UploadFileOptionFiles[]
  filePath?: string
  name?: string
  formData?: any
}

type UserBaseInfo = {
  id: number
  name: string
  score: number
  gender: string
  avatarUrl: string
}

type IUserInfo = {
  nickname?: string
  avatar?: string
  /** 微信的 openid，非微信没有这个字段 */
  openid?: string
  token?: string
}

enum TestEnum {
  A = 'a',
  B = 'b',
}
