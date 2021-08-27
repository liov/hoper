const process = require('child_process');

let gateway = process.execSync("go list -m -f {{.Dir}} github.com/grpc-ecosystem/grpc-gateway/v2",{
    cwd:'D:/hoper/server/go/lib'
}).toString()
console.log(gateway);