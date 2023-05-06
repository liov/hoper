import {API_HOST} from "@/env/config";
import type {MomentList} from "@/model/moment";
import uniHttp from "@/utils/request";
import {userStore} from "@/store";
import type {Response} from "@/service/response";

class MomentService {
    static async getMomentList(pageNo: number, pageSize: number): Promise<MomentList> {
        try {
            const {data} = await uniHttp.get<Response<MomentList>>("/api/v1/moment",
                {
                    pageNo,
                    pageSize
                },
                {
                    header: {
                        'custom-header': 'hello' //自定义请求头信息
                    }
                });
            console.log(data.details);
            if (data.details.users) userStore.appendUsers(data.details.users);
            return data.details;
        } catch (err: any) {
            await uni.showToast({
                title: decodeURI(err.errMsg),
                icon: "error",
                duration: 1000
            });
            console.error(err)
            return Promise.reject(err);
        }
    }
}
export default MomentService;