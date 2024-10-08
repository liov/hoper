
import 'dart:io';

var include = [];
Future<void> main() async{
  var protoPath = "../../proto";
  var goprojectPath = "../../server/go";
  const arguments = ["list", "-m","-f","{{.Dir}}"];
  var gopath = Platform.environment['GOPATH'];
/*  await Process.run("go",["mod" "download","github.com/googleapis/googleapis"],workingDirectory: goprojectPath);
 */
 var result = await Process.run("go",[...arguments,"github"
     ".com/hopeio/protobuf"],workingDirectory: goprojectPath);
 var cherryPath = '${(result.stdout as String).trimRight()}/_proto';

  include = ["-I$protoPath","-I$cherryPath"];
  Directory('${Directory.current.path}/lib/generated/protobuf').create();
  await generate(protoPath,[]);
  await generate(cherryPath,[]);
}

Future<void> generate(String dir,List<String> excludes) async {
  await for(var element in  Directory(dir).list()){
    if(await FileSystemEntity.type(element.path) == FileSystemEntityType.directory){
      if (excludes.any((e) => element.path.endsWith(e)))return;
      var  result = await Process.run("protoc",[...include,"${element.path}/*.proto","--dart_out=grpc:${Directory.current.path}/lib/generated/protobuf"]);
      print(result.stderr);
      print(result.stdout);
      await generate(element.path,[]);
    }
  }
}