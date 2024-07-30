import { API_HOST } from '@/env/config'
import type { MomentList } from '@/model/moment'
import uniHttp from '@/utils/request'
import { userStore } from '@/store'
import type { Response } from '@/service/response'
import moment from '@/pages/moment/moment_list.vue'

class UserService {
  static async active(id: number, secret: string): Promise<void> {
    await uniHttp.get<Response<void>>(`/api/v1/user/active/${id}/${secret}`)
  }
}
export default UserService
