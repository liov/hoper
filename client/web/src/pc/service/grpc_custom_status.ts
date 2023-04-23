import { GrpcStatusCode } from "@protobuf-ts/grpcweb-transport";

declare module "@protobuf-ts/grpcweb-transport" {
  export enum GrpcStatusCode {
    DBError = 21000,
  }
}

console.log("扩展grpc status");
const myEnumValues = GrpcStatusCode as any;
myEnumValues[21000] = "DBError";
