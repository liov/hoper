import { API_HOST } from '@/env/config'
import type { MomentList } from '@/model/moment'
import uniHttp from '@/utils/request'
import { userStore } from '@/store'
import type { Response } from '@/service/response'
import moment from '@/pages/moment/moment_list.vue'

class MomentService {
  static async getMomentList(pageNo: number, pageSize: number): Promise<MomentList> {
    const { data } = await uniHttp.get<Response<MomentList>>(
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
    console.log(data.details)
    if (data.details.users) userStore.appendUsers(data.details.users)
    if (data.details.list)
      data.details.list.forEach((moment) => {
        moment.user = userStore.getUser(moment.userId)
        if (moment.images && moment.images !== '') moment.imagesUrls = moment.images.split(',')
      })
    return data.details
  }
}
export default MomentService
