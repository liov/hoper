const process = require("child_process");
const fs = require("fs");
const path = require("path");

const projectpath = "D:/code/hoper/server/go/mod";

const goList = "go list -m -f {{.Dir}}";
process.execSync(`go mod download github.com/googleapis/googleapis`, {
  cwd: projectpath
});

function getDepPath(mod){
  return process.execSync(`${goList} ${mod}`, {
    cwd: projectpath
  }).toString().trimEnd()
}
const googleapis = getDepPath("github.com/googleapis/googleapis");
const libpath = getDepPath("github.com/liov/hoper/server/go/lib");
console.log(libpath);
const gateway = getDepPath("github.com/grpc-ecosystem/grpc-gateway/v2");
console.log(gateway);

const protopath = __dirname;
const libproto = libpath + "/protobuf";
const third = libpath + "/protobuf/third";

function include(isThird){
  if (isThird) return `-I${third}`;
  return  `-I${gateway} -I${googleapis} -I${protopath} -I${libproto} -I${third}`;
}

const dartConfig = {
  output: "D:/code/hoper\\client\\flutter\\lib\\generated\\protobuf",
  cwd: "D:/code/hoper\\client\\flutter",
  getCmd(filepath,isThird) {
   return `protoc ${include(isThird)} ${path.join(filepath, "*.proto")} --dart_out=grpc:${this.output}`
  },
};

function dartgenerate(dir, exclude, config,isThird=false) {
  fs.readdir(dir, function(err, files) {
    files.forEach(function(filename) {
      //获取当前文件的绝对路径
      let filepath = path.join(dir, filename);
      //根据文件路径获取文件信息，返回一个fs.Stats对象
      fs.stat(filepath, function(err2, stats) {
        if (err2) {
          console.warn(`获取文件stats失败,${filepath}`);
        } else {
          if (stats.isDirectory()) {
            if (exclude.includes(filename)) {
              return;
            }
            try {
              process.execSync(config.getCmd(filepath,isThird), {
                cwd: config.cwd
              });
            } catch (e) {
              console.log(e);
            }
            dartgenerate(filepath, [],config,isThird);
          }
        }
      });
    });
  });
}

dartgenerate(protopath, [], dartConfig);
dartgenerate(libproto, ["third"], dartConfig);
dartgenerate(third, [], dartConfig);
