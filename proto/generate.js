const childProcess = require("child_process");
const fs = require("fs");
const path = require("path");
const repo = "D:/code/hopeio/hoper"
const goProjectPath = repo+"/server/go";

const goList = "go list -m -f {{.Dir}}";

function getDepPath(mod){
  return childProcess.execSync(`${goList} ${mod}`, {
    cwd: goProjectPath
  }).toString().trimEnd()
}

const zetaPath = getDepPath("github.com/hopeio/zeta");
console.log(zetaPath);


const protopath = __dirname;
const zetaProto = zetaPath + "/protobuf/_proto";

const baseCmd = `protoc -I${protopath} -I${zetaProto}`

const goConfig = {
  output: repo+"\\server\\go\\protobuf",
  getCmd(filepath) {
    return [
       `${baseCmd} ${path.join(filepath, "*.proto")} --go-patch_out=plugin=go,paths=source_relative:${this.output}`,
      `${baseCmd} ${path.join(filepath, "*.service.proto")} --go-patch_out=plugin=go-grpc,paths=source_relative:${this.output}`,
      `${baseCmd} ${path.join(filepath, "*.enum.proto")} --enum_out=paths=source_relative:${this.output}`,
      `${baseCmd} ${path.join(filepath, "*.service.proto")} --grpc-gin_out=paths=source_relative:${this.output}`,
      `${baseCmd} ${path.join(filepath, "*.service.proto")} --openapiv2_out=logtostderr=true:${this.output}/api`,
      `${baseCmd} ${path.join(filepath, "*.service.proto")} --govalidators_out=paths=source_relative:${this.output}`,
      `${baseCmd} ${path.join(filepath, "*.service.proto")} --gql_out=svc=true,merge=true,paths=source_relative:${this.output}`,
      `${baseCmd} ${path.join(filepath, "*.service.proto")} --gogql_out=svc=true,merge=true,paths=source_relative:${this.output}`,
    ]
  },
};


const dartConfig = {
  output: repo+"\\client\\app\\lib\\generated\\protobuf",
  getCmd(filepath) {
   return `${baseCmd} ${path.join(filepath, "*.proto")} --dart_out=grpc:${this.output}`
  },
};

const grpcWebConfig = {
  output: repo+"\\client\\web\\generated\\grpc-web",
  getCmd(filepath) {
    return `${baseCmd}  ${path.join(filepath, "*.proto")} --js_out=import_style=commonjs,binary:${this.output} --grpc-web_out=import_style=typescript,mode=grpcwebtext:${this.output}`;
  },
};

const protobufTsConfig = {
  output: repo+"\\client\\web\\generated\\protobuf-ts",
  cwd: repo+"\\client\\web",
  getCmd(filepath) {
    return `npx ${baseCmd}  ${path.join(filepath, "*.proto")} --ts_out ${this.output}`;
  },
};

function generate(dir, exclude, config) {
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
              const cmd = config.getCmd(filepath);
              if (Array.isArray(cmd)) {
                cmd.forEach(function(subCmd){
                  childProcess.execSync(subCmd,{ cwd: config.cwd,encoding: 'utf-8' });
                })
              }else {
                childProcess.execSync(cmd,{ cwd: config.cwd,encoding: 'utf-8' });
              }
            } catch (e) {
              console.log(e.output.stderr);
            }
            generate(filepath, [],config);
          }
        }
      });
    });
  });
}

process.argv.slice(2).forEach(function(val,index,array){
  switch (val) {
    case 'go':
      //generate(protopath, [], goConfig);
      childProcess.execSync(`protogen.exe go -e -w -q -p ${protopath} -g ${goConfig.output}`,{ cwd: goConfig.output,encoding: 'utf-8' });
      break;
    case 'dart':
      generate(protopath, [], dartConfig);
      generate(zetaProto, [], dartConfig);
      break;
    case 'grpc-web':
      generate(protopath, [], grpcWebConfig);
      generate(zetaProto, [], grpcWebConfig);
      break;
    case 'protobuf-ts':
      generate(protopath, [], protobufTsConfig);
      generate(zetaProto, [], protobufTsConfig);
      break;
  }
})




