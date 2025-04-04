import { API_HOST } from '@/env/config'
import type { MomentList } from '@/model/moment'
import {unirequest} from '@hopeio/utils/uniapp'
import moment from '@/pages/moment/moment_list.vue'
import {ResData} from "@hopeio/utils/types";

class UserService {
  static async active(id: number, secret: string): Promise<void> {
    await unirequest.get<ResData<void>>(`/api/v1/user/active/${id}/${secret}`)
  }
}
export default UserService
