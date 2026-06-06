
import { httpclient } from '@/api/common'
import { MomentListResp } from '@gen/pb/content/moment.service'
import { useUserStore } from '@/store/user'
import { CommonResp } from '@hopeio/utils/types'

const userStore = useUserStore()
class MomentService {
  static async getMomentList(pageNo: number, pageSize: number): Promise<CommonResp<MomentListResp>> {
    const data = await httpclient.get<CommonResp<MomentListResp>>(
      '/api/moment',
      { query: { pageNo, pageSize }, decode: MomentListResp },
    )
    console.log(data.data)
    if (data.data?.users) userStore.appendUsers(data.data.users)
    if (data.data?.list)
      data.data.list.forEach((moment) => {
        moment.user = userStore.getUser(moment.userId)
      })
    return data
  }
}
export default MomentService
