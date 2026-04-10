import { API_HOST } from '@/env/config'
import type { MomentList } from '@/model/moment'
import {httpclient} from '@hopeio/utils/uniapp'

import moment from '@/pages/moment/moment_list.vue'
import { useUserStore } from '@/store/user'
import { CommonResp } from '@hopeio/utils/types'

const userStore = useUserStore()
class MomentService {
  static async getMomentList(pageNo: number, pageSize: number): Promise<MomentList> {
    const { data } = await httpclient.get<CommonResp<MomentList>>(
      '/api/moment',
      { query: { pageNo, pageSize } },
    )
    console.log(data)
    if (data.users) userStore.appendUsers(data.users)
    if (data.list)
      data.list.forEach((moment) => {
        moment.user = userStore.getUser(moment.userId)
        if (moment.images && moment.images !== '') moment.imagesUrls = moment.images.split(',')
      })
    return data
  }
}
export default MomentService
