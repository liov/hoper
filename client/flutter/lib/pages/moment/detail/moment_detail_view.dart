import 'package:app/components/async/async.dart';
import 'package:app/service/moment.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:app/generated/protobuf/content/content.model.pb.dart' as $pb;
import 'moment_detail_controller.dart';
import 'package:fixnum/fixnum.dart';
import 'package:app/generated/protobuf/utils/request/param.pb.dart' as $1;

class MomentDetailView extends StatelessWidget {

  MomentDetailView({this.moment,required this.id}):super();

  final MomentClient momentClient = Get.find();
  final $pb.Moment? moment;
  final Int64 id;

  Future<$pb.Moment> getMoment() async{
     if(moment!=null) return moment!;
   final rpcMoment =  await momentClient.stub.info($1.Object(id:id));
   return rpcMoment;
  }

  @override
  Widget build(BuildContext context) {
    final future = getMoment();
    return FutureBuilder<$pb.Moment>(
      future: future,
      builder:(BuildContext context, AsyncSnapshot<$pb.Moment> snapshot) {
        final noReady = snapshot.handle();
        if (noReady != null) return Scaffold(body:noReady);
        final moment = snapshot.data!;
        return  Scaffold(body:Center(child: Text(moment.content),));
      });
  }
}

