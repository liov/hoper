
import {HelloRequest,HelloReply} from './protobuf/helloworld_pb';
import {GreeterClient} from './protobuf/helloworld_grpc_pb';
import * as grpc from "grpc";

function main() {
    const client = new GreeterClient('localhost:50051',  grpc.credentials.createInsecure());
    const request = new HelloRequest();
    request.setName('web client!');
    client.sayHello(request, (err:any,response:HelloReply|undefined) => {
        console.log('Greeting:', response? response.getMessage():err);
    });
}

main();