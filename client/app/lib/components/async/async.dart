
import 'package:flutter/material.dart';

import '../../global/service.dart';

extension HandleAsync<T> on AsyncSnapshot<T> {
  Widget? handle<T>(){
    switch (connectionState) {
      case ConnectionState.none:
      case ConnectionState.active:
        return Text('ConnectionState.active');
      case ConnectionState.waiting:
        return Center(
          child: CircularProgressIndicator(),
        );
      case ConnectionState.done:
        globalService.logger.d('done');
        return null;
    }
  }
}