
import 'package:flutter/material.dart';

extension HandleAsync<T> on AsyncSnapshot<T> {
  Widget? handle<T>(){
    switch (this.connectionState) {
      case ConnectionState.none:
      case ConnectionState.active:
        return Text('ConnectionState.active');
      case ConnectionState.waiting:
        return Center(
          child: CircularProgressIndicator(),
        );
      case ConnectionState.done:
        return null;
    }
  }
}