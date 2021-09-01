const process = require('child_process');
const fs = require('fs');
const path = require('path');

const projectpath = 'D:/code/hoper/server/go/mod'

const goList = "go list -m -f {{.Dir}}"
process.execSync(`go mod download github.com/googleapis/googleapis`, {
    cwd: projectpath
})
const googleapis = process.execSync(`${goList} github.com/googleapis/googleapis`, {
    cwd: projectpath
}).toString().trimEnd()
const libpath = process.execSync(`${goList} github.com/liov/hoper/server/go/lib`, {
    cwd: projectpath
}).toString().trimEnd()
console.log(libpath)
const gateway = process.execSync(`${goList} github.com/grpc-ecosystem/grpc-gateway/v2`, {
    cwd: libpath
}).toString().trimEnd()
console.log(gateway)

const protobuf = process.execSync(`${goList} google.golang.org/protobuf`, {
    cwd: libpath
}).toString().trimEnd()
const protopath = __dirname
const libproto = libpath + "/protobuf"
const third = libpath + "/protobuf/third"
let cmd = `protoc -I${gateway} -I${googleapis} -I${protobuf} -I${protopath} -I${libproto} -I${third} `

function dartgenerate(dir,exlude) {
    fs.readdir(dir, function (err, files) {
        files.forEach(function (filename) {
            //获取当前文件的绝对路径
            let filepath = path.join(dir, filename);
            //根据文件路径获取文件信息，返回一个fs.Stats对象
            fs.stat(filepath, function (err2, stats) {
                if (err2) {
                    console.warn(`获取文件stats失败,${filepath}`);
                } else {
                    if (stats.isDirectory()) {
                        if(exlude.includes(filename)){
                            return
                        }
                        try {
                            process.execSync(`${cmd} ${path.join(filepath, '*.proto')} --dart_out=grpc:D:/code/hoper\\client\\flutter\\lib\\generated\\protobuf`, {
                                cwd: "D:/code/hoper\\client\\flutter"
                            })
                        } catch (e) {
                        }
                        dartgenerate(filepath,[])
                    }
                }
            })
        });
    })
}

dartgenerate(protopath,[])
dartgenerate(libpath + "/protobuf",["third","utils"])
cmd = `protoc -I${gateway}  -I${protobuf} -I${protopath} -I${third} `
dartgenerate(libpath + "/protobuf/third",[])
