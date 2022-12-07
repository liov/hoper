
import 'dart:io';

var include = [];
Future<void> main() async{
  var protoPath = "../../proto";
  var goprojectPath = "../../server/go/mod";
  const arguments = ["list", "-m","-f","{{.Dir}}"];
  const googleapis = "github.com/googleapis/googleapis@v0.0.0-20220520010701-4c6f5836a32f";
  var gopath = Platform.environment['GOPATH'];
/*  await Process.run("go",["mod" "download","github.com/googleapis/googleapis"],workingDirectory: goprojectPath);
 */
 var result = await Process.run("go",[...arguments,"github.com/liov/hoper/server/go/lib"],workingDirectory: goprojectPath);
  var golibPath = (result.stdout as String).trimRight();
  var google = gopath! + "pkg/mod/" + googleapis;
  var googleexists = await Directory(google).exists();
  if (!googleexists) {
    await Process.run("go",["get",googleapis],workingDirectory: goprojectPath);
    result = await Process.run("go",[...arguments,"github.com/googleapis/googleapis"],workingDirectory: goprojectPath);
    google = (result.stdout as String).trimRight();
    await Process.run("go",["mod","tidy"],workingDirectory: goprojectPath);
  }

  result = await Process.run("go",[...arguments,"github.com/grpc-ecosystem/grpc-gateway/v2"],workingDirectory: golibPath);
  var gateway = (result.stdout as String).trimRight();
  include = ["-I${google}","-I${gateway}","-I${protoPath}","-I${golibPath}/protobuf","-I${golibPath}/protobuf/third"];
  Directory('${Directory.current.path}/lib/generated/protobuf').create();
  await generate(protoPath,[]);
  await generate(golibPath+"/protobuf",["third","utils"]);
  include = ["-I${gateway}","-I${protoPath}","-I${golibPath}/protobuf/third"];
  await generate(golibPath+"/protobuf/third",[]);
}

Future<void> generate(String dir,List<String> exludes) async {
  await for(var element in  Directory(dir).list()){
    if(await FileSystemEntity.type(element.path) == FileSystemEntityType.directory){
      if (exludes.any((e) => element.path.endsWith(e)))return;
      var  result = await Process.run("protoc",[...include,"${element.path}/*.proto","--dart_out=grpc:${Directory.current.path}/lib/generated/protobuf"]);
      print(result.stderr);
      print(result.stdout);
      await generate(element.path,[]);
    }
  }
}