
import 'package:app/global/service.dart';
import 'package:app/model/response.dart';
import 'package:dio/dio.dart';
import 'package:applib/util/json.dart';
import '../model/weibo/weibo.dart';

const prefix = "https://weibo.com/ajax/statuses";

class WeiboClient {

    WeiboClient(){
     httpClient.interceptors.add(
        InterceptorsWrapper(
          onRequest: (RequestOptions options, RequestInterceptorHandler handler) {
            return handler.next(options);
          },
          onResponse: (Response response, ResponseInterceptorHandler handler) {
            return handler.next(response);
          },
          onError: (DioException error, ErrorInterceptorHandler handler) {
            globalService.logger.e(error);
            return handler.next(error);
          },
        ),
      );
    }
    final Dio httpClient = Dio(BaseOptions(
        baseUrl: prefix,
        connectTimeout: const Duration(seconds: 5),
        receiveTimeout: const Duration(seconds: 3),
        headers: {
          "Cookie":"SCF=Am5M98SU-7_OGKs36KFt79uErGYQl85tXBdbV97xFuP5AAreLz_0QBgiRxarBfwC2Krv8MuwQf7qICN33i_n7vo.; SINAGLOBAL=68123827669.85148.1747380261721; ULV=1747380261748:1:1:1:68123827669.85148.1747380261721:; XSRF-TOKEN=p-VRsnJjVjKRghXhyzSRBNWS; ALF=1751533488; SUB=_2A25FOsrgDeRhGeNM6FIY8i3LzTyIHXVmNkIorDV8PUJbkNAbLWb1kW1NTiTnsQ91IB1dhEIdBJ9ZC48WKFFWnxF2; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WWoVxvm6ETYDZJXPaqQd5i35JpX5KMhUgL.Fo-Ee054eoeNSo52dJLoIEBLxKnLBoBLB.-LxKBLBonL1h5LxK.LBK-LB.eLxK-L1K5L1Kqt; WBPSESS=SwpthAMfiMQK-4Y9U-w8UGm6JrmGzTyRUros7m58pBjoEdWYaTNh-RHfXpOLUvbKxa22CnfNXUDMEhneHTMcXs7F8BisQZmzPsrGw_Zb1qbMykwM2obeKb-8f_gjob5YeAKS819IQXi1ne1DCB3IIw==",
          "Content-Type": 'application/json',
          "Connection": 'keep-alive',
        },
    ));


    Future<WeiboList?> getList({required int uid,required int page,required int feature,String? sinceId}) async {
      globalService.logger.d("uid:$uid,page:$page,feature:$feature,sinceId:$sinceId");
    final response = await httpClient.get('/mymblog',queryParameters: {
      'uid':uid,
      'page':page,
      'feature':feature,
      'since_id':sinceId
    });
    return response.getWeiboData((v)=>WeiboList.fromJson(v  as Map<String,dynamic>));
  }
  Future<WeiboOriginList?> getOriginalList({required int uid,required int page,required int feature,String? sinceId}) async {
    final response = await httpClient.get('/searchProfile',queryParameters: {
      'uid':uid,
      'page':page,
      'hasori':1
    });
    return response.getWeiboData((v)=>WeiboOriginList.fromJson(v  as Map<String,dynamic>));
  }
}
