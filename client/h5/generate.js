const process = require("child_process");
const fs = require("fs");
const path = require("path");

const protopath = "D:/code/hoper/proto";

const libpath = "D:/code/hoper/server/go/lib";
const libproto = libpath + "/protobuf";
const third = libpath + "/protobuf/third";

function include(isThird) {
  if (isThird) return `-I${third}`;
  return `-I${protopath} -I${libproto} -I${third}`;
}

const grpcWebConfig = {
  output: "D:/code/hoper\\client\\h5\\generated\\grpc-web",
  cwd: "D:/code/hoper\\client\\h5",
  getCmd(filepath, isThird) {
    return `protoc ${include(isThird)}  ${path.join(
      filepath,
      "*.proto"
    )} --js_out=import_style=commonjs,binary:${
      this.output
    } --grpc-web_out=import_style=typescript,mode=grpcwebtext:${this.output}`;
  },
};

const protobufTsConfig = {
  output: "D:/code/hoper\\client\\h5\\generated\\protobuf-ts",
  cwd: "D:/code/hoper\\client\\h5",
  getCmd(filepath, isThird) {
    return `npx protoc ${include(isThird)}  ${path.join(
      filepath,
      "*.proto"
    )} --ts_out ${this.output}`;
  },
};

function generate(dir, exclude, config, isThird = false) {
  fs.readdir(dir, function (err, files) {
    files.forEach(function (filename) {
      //获取当前文件的绝对路径
      const filepath = path.join(dir, filename);
      //根据文件路径获取文件信息，返回一个fs.Stats对象
      fs.stat(filepath, function (err2, stats) {
        if (err2) {
          console.warn(`获取文件stats失败,${filepath}`);
        } else {
          if (stats.isDirectory()) {
            if (exclude.includes(filename)) {
              return;
            }
            try {
              process.execSync(config.getCmd(filepath, isThird), {
                cwd: config.cwd,
              });
            } catch (e) {
              console.log(e);
            }
            generate(filepath, [], config, isThird);
          }
        }
      });
    });
  });
}

generate(protopath, [], protobufTsConfig);
generate(libproto, ["third"], protobufTsConfig);
generate(third, [], protobufTsConfig, true);
