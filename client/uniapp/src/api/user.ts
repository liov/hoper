

import { httpclient } from '@/api/common'

import {CommonResp} from "@hopeio/utils/types";

class UserService {
  static async active(id: number, secret: string): Promise<void> {
    await httpclient.get<CommonResp<void>>(`/api/user/active/${id}/${secret}`)
  }
  static async login(params: any): Promise<CommonResp<any>> {
    return await httpclient.post<CommonResp<any>>('/api/user/login', { data: params })
  }
  static async signup(params: any): Promise<CommonResp<any>> {
    return await httpclient.post<CommonResp<any>>('/api/user', { data: params })
  }
  static async baseUserList(ids: number[]): Promise<CommonResp<any>> {
    return await httpclient.post<CommonResp<any>>('/api/user/baseUserList', { data: { ids } })
  }
}
export default UserService
