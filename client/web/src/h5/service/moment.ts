import { MomentServiceClient } from "@generated/protobuf-ts/content/moment.service.client";
import type { MomentListRep } from "@generated/protobuf-ts/content/moment.service";
import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport";
import { showFailToast } from "vant";
import type { RpcError } from "@protobuf-ts/runtime-rpc";

const momentClient = new MomentServiceClient(
  new GrpcWebFetchTransport({
    baseUrl: "https://grpc.hoper.xyz",
  })
);

export async function momentList(
  pageNo: number,
  pageSize: number
): Promise<MomentListRep> {
  try {
    const { response, status } = await momentClient.list({ pageNo, pageSize });
    console.log(status);
    return response;
  } catch (e) {
    const rpcError = e as RpcError;
    showFailToast(decodeURI(rpcError.message));
    return Promise.reject(rpcError);
  }
}
