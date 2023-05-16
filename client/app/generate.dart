
import 'dart:io';

var include = [];
Future<void> main() async{
  var protoPath = "../../proto";
  var goprojectPath = "../../server/go";
  const arguments = ["list", "-m","-f","{{.Dir}}"];
  var gopath = Platform.environment['GOPATH'];
/*  await Process.run("go",["mod" "download","github.com/googleapis/googleapis"],workingDirectory: goprojectPath);
 */
 var result = await Process.run("go",[...arguments,"github.com/hopeio/pandora"],workingDirectory: goprojectPath);
 var pandoraPath = (result.stdout as String).trimRight()+'/protobuf/_proto';

  include = ["-I${protoPath}","-I${pandoraPath}"];
  Directory('${Directory.current.path}/lib/generated/protobuf').create();
  await generate(protoPath,[]);
  await generate(pandoraPath,[]);
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