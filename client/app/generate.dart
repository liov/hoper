import 'dart:io';

var include = [];

Future<void> main() async {
  var protoPath = "../../proto";
  var goprojectPath = "../../server/go";
  const arguments = ["list", "-m", "-f", "{{.Dir}}"];


  // 注意：原代码中 go mod download 被注释了，如果需要请取消注释并修复语法错误 (原代码少了一个逗号)
  // await Process.run("go", ["mod", "download", "github.com/googleapis/googleapis"], workingDirectory: goprojectPath);

  var result = await Process.run("go", [...arguments, "github.com/hopeio/protobuf"], workingDirectory: goprojectPath);

  if (result.exitCode != 0) {
    print("Error getting go module path: ${result.stderr}");
    exit(1);
  }

  var cherryPath = '${(result.stdout as String).trim()}/_proto';

  // 规范化路径，避免相对路径混乱
  include = [
    "-I${protoPath}",
    "-I${Directory(cherryPath).absolute.path}"
  ];

  var outputDir = Directory('${Directory.current.path}/lib/generated/protobuf');
  if (!await outputDir.exists()) {
    await outputDir.create(recursive: true);
  }

  // 分别生成
  await generate(protoPath, []);
  await generate(cherryPath, []);
}

Future<void> generate(String dir, List<String> excludes) async {
  var dirEntity = Directory(dir);
  if (!await dirEntity.exists()) {
    print("Warning: Directory $dir does not exist, skipping.");
    return;
  }

  // 使用 list(recursive: true) 一次性获取所有文件，或者保持你的递归逻辑但只处理文件
  // 这里优化逻辑：直接查找当前目录下的 .proto 文件，不再依赖 shell 通配符
  await for (var entity in dirEntity.list()) {
    if (entity is Directory) {
      if (excludes.any((e) => entity.path.endsWith(e))) continue;

      // 递归处理子目录
      await generate(entity.path, excludes);
    } else if (entity is File && entity.path.endsWith('.proto')) {
      // 找到具体的 .proto 文件
      print("Compiling: ${entity.path}");

      // 构造命令参数：明确传入文件路径，不使用 *.proto

      var result = await Process.run("protoc", [
        ...include,
        "--dart_out=grpc:${Directory.current.path}/lib/generated/protobuf",
        entity.path // 直接传入具体文件
      ]);

      if (result.exitCode != 0) {
        print("Error compiling ${entity.path}:");
        print(result.stderr);
        print(result.stdout);
      } else {
        print("Success: ${entity.path}");
      }
    }
  }
}