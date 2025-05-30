
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
          "Cookie":"XSRF-TOKEN=Guf7X9miwgJnZWSb8XGl4aoH; SUB=_2AkMTABV5f8NxqwJRmfwVzGzgZY1-wgnEieKlXOSiJRMxHRl-yT9vqmk5tRB6OIA7lbl9FTWqfOwFzW72UDoJ7Im94AtD; SUBP=0033WrSXqPxfM72-Ws9jqgMF55529P9D9WWclODofmRcG7uVnHn2aqpb; WBPSESS=2jFn3n4I-3CwFbURoTaQu2CrqYdzybdxRQfZIAoEvvrNsMFxjpAkHDJpo96w63LIyO1qEtcZbPflvFottVCHmmQsJvkH8o5zmEDAykoR9hU2OmCtZLtYf-MjgyyS92ddJ75iujgZQd3jj-yjqEwl1Pht5kw2k1-MDB2pm3_FcwI=",
          "Content-Type": 'application/json',
          "Connection": 'keep-alive',
        },
    ));


    Future<Map<String,dynamic>?> getList({required int uid,required int page,required int feature,String? sinceId}) async {
      globalService.logger.d("uid:$uid,page:$page,feature:$feature,sinceId:$sinceId");
    final response = await httpClient.get('/mymblog',queryParameters: {
      'uid':uid,
      'page':page,
      'feature':feature,
      'since_id':sinceId
    });
    return response.getWeiboData((v)=>v  as Map<String,dynamic>);
  }
}
