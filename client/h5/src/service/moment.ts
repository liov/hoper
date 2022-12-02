import { MomentServiceClient } from "@generated/protobuf-ts/content/moment.service.client";
import type { MomentListRep } from "@generated/protobuf-ts/content/moment.service";
import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport";
import type { RpcError } from "grpc-web";
import { Toast } from "vant";

const momentClient = new MomentServiceClient(
  new GrpcWebFetchTransport({
    baseUrl: "http://localhost:8090",
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
  } catch (e: RpcError) {
    Toast.fail(decodeURI(e.message));
    return Promise.resolve(e);
  }
}
