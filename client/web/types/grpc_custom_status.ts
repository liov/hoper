import { GrpcStatusCode } from "@protobuf-ts/grpcweb-transport";

declare module "@protobuf-ts/grpcweb-transport" {
  export enum GrpcStatusCode {
    DBError = 21000,
  }
}

console.log("扩展grpc status");
const myEnumValues = GrpcStatusCode as any;
myEnumValues[GrpcStatusCode.DBError] = "DBError";
