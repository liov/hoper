import {API_HOST} from '@/env/config'
import type {MomentList} from '@/model/moment'
import {UniRequest} from '@hopeio/utils/uniapp'
import moment from '@/pages/moment/moment_list.vue'
import {ResData} from "@hopeio/utils/types";

class WeiboService {
  private static unirequest = new UniRequest({
    baseUrl: "https://weibo.com",
    withCredentials: true,
    header: {
      'content-type': 'application/json',
      Cookie: "SCF=Am5M98SU-7_OGKs36KFt79uErGYQl85tXBdbV97xFuP5AAreLz_0QBgiRxarBfwC2Krv8MuwQf7qICN33i_n7vo.; SINAGLOBAL=68123827669.85148.1747380261721; ULV=1747380261748:1:1:1:68123827669.85148.1747380261721:; XSRF-TOKEN=p-VRsnJjVjKRghXhyzSRBNWS; ALF=1751533488; SUB=_2A25FOsrgDeRhGeNM6FIY8i3LzTyIHXVmNkIorDV8PUJbkNAbLWb1kW1NTiTnsQ91IB1dhEIdBJ9ZC48WKFFWnxF2; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WWoVxvm6ETYDZJXPaqQd5i35JpX5KMhUgL.Fo-Ee054eoeNSo52dJLoIEBLxKnLBoBLB.-LxKBLBonL1h5LxK.LBK-LB.eLxK-L1K5L1Kqt; WBPSESS=SwpthAMfiMQK-4Y9U-w8UGm6JrmGzTyRUros7m58pBjoEdWYaTNh-RHfXpOLUvbKxa22CnfNXUDMEhneHTMcXs7F8BisQZmzPsrGw_Zb1qbMykwM2obeKb-8f_gjob5YeAKS819IQXi1ne1DCB3IIw=="
      //'user-agent': 'uniapp-' + uni.getAppBaseInfo().appName,
    },
    timeout: 10000,
  })

  static async list(uid: number, page: number): Promise<any> {
    await this.unirequest.get<ResData<void>>(`/ajax/statuses/searchProfile`, {
      query: {
        uid,
        page,
        hasori: 1
      }
    })
  }
}

export default WeiboService
