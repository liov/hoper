

import { httpclient } from '@/api/common'
import { LoginReq, LoginResp, SignupReq, SingUpVerifyReq } from '@gen/pb/user/user.service';

import {CommonResp} from "@hopeio/utils/types";

class UserService {
  static async active(id: number, secret: string): Promise<void> {
    await httpclient.get<CommonResp<void>>(`/api/user/active/${id}/${secret}`)
  }
  static async login(params: LoginReq): Promise<LoginResp> {
    return await httpclient.post<LoginResp>('/api/user/login', { data: params })
  }
  /** 注册前校验邮箱/手机号是否可用 */
  static async signupVerify(params: SingUpVerifyReq): Promise<CommonResp<string>> {
    return await httpclient.post<CommonResp<string>>('/api/user/signupVerify', { data: params })
  }
  /** 发送验证码：服务端根据 mail / phone 区分下发渠道；action：1=注册 */
  static async sendVerifyCode(params: { mail?: string; phone?: string; action: number; vCode: string }): Promise<CommonResp<void>> {
    return await httpclient.get<CommonResp<void>>('/api/sendVerifyCode', { query: params })
  }
  static async signup(params: SignupReq): Promise<CommonResp<string>> {
    return await httpclient.post<CommonResp<string>>('/api/user', { data: params })
  }
  static async baseUserList(ids: number[]): Promise<CommonResp<any>> {
    return await httpclient.post<CommonResp<any>>('/api/user/baseUserList', { data: { ids } })
  }
}
export default UserService
