var PROTO_PATH = __dirname + '/../../../protobuf/helloworld.proto';

var grpc = require('grpc');
var protoLoader = require('@grpc/proto-loader');
var packageDefinition = protoLoader.loadSync(
    PROTO_PATH,
    {keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true
    });
var hello_proto = grpc.loadPackageDefinition(packageDefinition).protobuf;

/**
 * Implements the SayHello RPC method.
 */
function sayHello(call, callback) {
    callback(null, {message: 'node ' + call.request.name});
}

/**
 * Starts an RPC server that receives requests for the Greeter service at the
 * sample server port
 */
function main() {
    var server = new grpc.Server();
    server.addService(hello_proto.Greeter.service, {sayHello: sayHello});
    server.bind('0.0.0.0:50051', grpc.ServerCredentials.createInsecure());
    server.start();
}

main();