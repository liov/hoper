using System;
using System.Threading.Tasks;
using Grpc.Core;
using Protobuf;

namespace hello
{
    public class GreeterService : Greeter.GreeterBase
    {
        public override Task<HelloReply> SayHello(HelloRequest request, ServerCallContext context)
        {
            return Task.FromResult(new HelloReply
            {
                Message = "Csharp " + request.Name
            });
        }
    }
}
