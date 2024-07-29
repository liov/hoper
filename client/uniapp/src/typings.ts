// 全局要用的类型放到这里

type IResData<T> = {
  code: number
  msg: string
  data: T
}

// uni.uploadFile文件上传参数
type IUniUploadFileOptions = {
  file?: File
  files?: UniApp.UploadFileOptionFiles[]
  filePath?: string
  name?: string
  formData?: any
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

declare interface Window {
  // window对象属性
  turnstile: any // 加入对象

  wx: any
  WeixinJSBridge: any
  __wxjs_environment: string
}
