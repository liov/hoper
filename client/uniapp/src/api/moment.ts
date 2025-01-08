import { API_HOST } from '@/env/config'
import type { MomentList } from '@/model/moment'
import {unirequest} from '@hopeio/utils/uniapp'

import moment from '@/pages/moment/moment_list.vue'
import { useUserStore } from '@/store/user'
const userStore = useUserStore()
class MomentService {
  static async getMomentList(pageNo: number, pageSize: number): Promise<MomentList> {
    const { data } = await unirequest.get<ResData<MomentList>>(
      '/api/v1/moment',
      {
        pageNo,
        pageSize,
      },
      {
        header: {
          'custom-header': 'hello', // 自定义请求头信息
        },
      },
    )
    console.log(data.data)
    if (data.data.users) userStore.appendUsers(data.data.users)
    if (data.data.list)
      data.data.list.forEach((moment) => {
        moment.user = userStore.getUser(moment.userId)
        if (moment.images && moment.images !== '') moment.imagesUrls = moment.images.split(',')
      })
    return data.data
  }
}
export default MomentService
