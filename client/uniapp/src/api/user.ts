

import { httpclient } from '@/api/common'
import { StringValue } from '@gen/pb/google/protobuf/wrappers';
import { Auth } from '@gen/pb/user/user.model';
import { LoginReq, LoginResp, SignupReq, SingUpVerifyReq } from '@gen/pb/user/user.service';

import {CommonResp} from "@hopeio/utils/types";

class UserService {
  static async active(id: number, secret: string): Promise<void> {
    await httpclient.get<CommonResp<void>>(`/api/user/active/${id}/${secret}`)
  }
  static async login(params: LoginReq): Promise<LoginResp> {
    return await httpclient.post<LoginResp>('/api/user/login', { data: params,decode: LoginResp })
  }
  /** 注册前校验邮箱/手机号是否可用 */
  static async signupVerify(params: SingUpVerifyReq): Promise<void> {
    return await httpclient.post<void>('/api/user/signupVerify', { data: params })
  }
  /** 发送验证码：服务端根据 mail / phone 区分下发渠道；action：1=注册 */
  static async sendVerifyCode(params: { mail?: string; phone?: string; action: number; vCode: string }): Promise<void> {
    return await httpclient.get<void>('/api/sendVerifyCode', { query: params })
  }
  static async signup(params: SignupReq): Promise<StringValue> {
    return await httpclient.post<StringValue>('/api/user', { data: params,decode: StringValue })
  }
  static async baseUserList(ids: number[]): Promise<any> {
    return await httpclient.post<any>('/api/user/baseUserList', { data: { ids } })
  }
  static async auth(): Promise<Auth> {
    return await httpclient.get<Auth>('/api/auth', { decode: Auth})
  }
}
export default UserService
