import { httpclient } from "./common";
import { showFailToast } from "vant";
import type { CommonResp } from "@hopeio/utils/types";
import type { MomentListResp } from "@gen/pb/content/moment.service";

export async function momentList(
  pageNo: number,
  pageSize: number
){
  return await httpclient.get<CommonResp<MomentListResp>>(`/api/moment?pageNo=${pageNo}&pageSize=${pageSize}`);
}
