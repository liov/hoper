import * as grpcWeb from 'grpc-web';
import {HelloRequest,HelloReply} from './protobuf/helloworld_pb';
import {GreeterClient} from './protobuf/HelloworldServiceClientPb';

const client = new GreeterClient('http://172.27.175.186:8080', null, null);

const request = new HelloRequest();
request.setName('web client!');

const call = client.sayHello(request, {'custom-header-1': 'value1'},
    (err: grpcWeb.Error, response: HelloReply) => {
        console.log(response.getMessage());
    });
call.on('status', (status: grpcWeb.Status) => {
    // ...
});