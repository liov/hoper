
import 'dart:io';
import 'package:crypto/crypto.dart';


import 'package:app/model/upload.dart';
import 'package:app/model/response.dart';
import 'package:app/utils/dio.dart';
import 'package:dio/dio.dart';

import 'package:grpc/src/client/call.dart';

import 'package:app/utils/observer.dart';

class UploadClient extends Observer<CallOptions> {

  UploadClient(Subject<CallOptions> subject){
    update(subject.options);
    subject.attach(this);
  }


  Future<String> getMD5(File file) async {
    final mdFive = await md5.bind(file.openRead()).first;
    return mdFive.toString();
  }

  Future<String> exists(String mdFive,int size) async {
    var api = '/v1/exists/${mdFive}/${size}';
    try {
      final response = await httpClient.get(api);
      return response.data;
    } catch (e) {
      print(e);
      return '';
    }
  }

  Future<String> upload(File file) async {
    final size = await file.length();
    final mdFive = await md5.bind(file.openRead()).first;
    final url = await exists(mdFive.toString(),size);
    if (url!='') return url;
    var api = '/v1/upload/${mdFive}';
    try {
      final response = await httpClient.post(api,
          data: FormData.fromMap({'file': await MultipartFile.fromFile(file.path)}),
       );
      return UploadInfo.fromJson(response.getData()).url;
    } catch (e) {
      print(e);
      return '';
    }
  }

  Future<String> uploadMultiple(List<File> files) async {
    final sizeList = await Future.wait(files.map((f) => f.length()));
    final md5List = await Future.wait(files.map((f) => md5.bind(f.openRead()).first.then((d) => d.toString())));
    final fileList = await Future.wait(files.map((f) =>MultipartFile.fromFile(f.path)));
    var api = '/v1/multiUpload';
    try {
      final response = await httpClient.post(api,
        data: FormData.fromMap({'size[]':sizeList,'md5[]':md5List,'file[]':fileList}),
         );
      return UploadInfo.fromJson(response.getData()).url;
    } catch (exception) {
      return '';
    }
  }

  @override
  void update(CallOptions? options) {
    // TODO: implement update
  }
}