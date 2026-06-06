
import 'package:flutter/material.dart';

extension HandleAsync<T> on AsyncSnapshot<T> {
  Widget? handle<T>({Widget? child}){
    switch (connectionState) {
      case ConnectionState.none:
      case ConnectionState.active:
        return Text('ConnectionState.active');
      case ConnectionState.waiting:
        return Center(
          child: child ?? CircularProgressIndicator(),
        );
      case ConnectionState.done:
        return null;
    }
  }
}