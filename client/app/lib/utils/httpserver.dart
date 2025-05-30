import 'dart:io';
import 'dart:async';
import 'dart:typed_data';

import 'package:app/utils/http_static/static_handler.dart';
import 'package:flutter/services.dart' show rootBundle;


import 'package:shelf/src/request.dart';
import 'package:shelf/src/response.dart';
import 'package:shelf/src/handler.dart';
import 'package:shelf_proxy/shelf_proxy.dart';


///This class allows you to create a simple server on `http://localhost:[port]/` in order to be able to load your assets file on a server. The default [port] value is `8080`.
class LocalhostServer {
  bool _started = false;
  HttpServer? _server;
  int _port = 8080;
  String _path = "assets/dist";

  LocalhostServer({int port = 8080, String path = "assets/dist"}) {
    _port = port;
    _path = path;
  }

  ///Starts the server on `http://localhost:[port]/`.
  ///
  ///**NOTE for iOS**: For the iOS Platform, you need to add the `NSAllowsLocalNetworking` key with `true` in the `Info.plist` file (See [ATS Configuration Basics](https://developer.apple.com/library/archive/documentation/General/Reference/InfoPlistKeyReference/Articles/CocoaKeys.html#//apple_ref/doc/uid/TP40009251-SW35)):
  ///```xml
  ///<key>NSAppTransportSecurity</key>
  ///<dict>
  ///    <key>NSAllowsLocalNetworking</key>
  ///    <true/>
  ///</dict>
  ///```
  ///The `NSAllowsLocalNetworking` key is available since **iOS 10**.
  Future<void> start() async {
    if (_started) {
      throw Exception('Server already started on http://localhost:$_port');
    }
    _started = true;

    var completer = Completer();

    runZonedGuarded(() {
      HttpServer.bind('127.0.0.1', _port).then((server) {
        print('Server running on http://localhost:' + _port.toString());

        _server = server;

        server.listen((HttpRequest request) async {
          Uint8List body = Uint8List(0);

          var path = _path + request.requestedUri.path;
          path = (path.startsWith('/')) ? path.substring(1) : path;
          path += (path.endsWith('/')) ? 'index.html' : '';

          try {
            body = (await rootBundle.load(path)).buffer.asUint8List();
          } catch (e) {
            print(e.toString());
            request.response.close();
            return;
          }

          var contentType = ['text', 'html'];
          if (!request.requestedUri.path.endsWith('/') &&
              request.requestedUri.pathSegments.isNotEmpty) {
            var mimeType = DefaultMimeTypeResolver.lookup(request.requestedUri.path);
            if (mimeType != null) {
              contentType = mimeType.split('/');
            }
          }

          request.response.headers.contentType =
              ContentType(contentType[0], contentType[1], charset: 'utf-8');
          request.response.add(body);
          request.response.close();
        });

        completer.complete();
      });
    }, (e, stackTrace) => print('Error: $e $stackTrace'));

    return completer.future;
  }

  ///Closes the server.
  Future<void> close() async {
    if (_server == null) {
      return;
    }
    await _server!.close(force: true);
    print('Server running on http://localhost:$_port closed');
    _started = false;
    _server = null;
  }

  ///Indicates if the server is running or not.
  bool isRunning() {
    return _server != null;
  }
}

Handler apiHandler(String proxyName, {prefix = "/api", dist = "assets/dist"}) {
  final ph = proxyHandler(proxyName);
  final ch = createStaticHandler(dist, defaultDocument: 'index.html',assets:true);
  return (Request request) {
    if (request.requestedUri.path.startsWith(prefix)) {
      return ph(request);
    }
    return ch(request);
  };
}

