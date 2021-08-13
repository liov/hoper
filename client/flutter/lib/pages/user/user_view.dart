import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:getwidget/components/list_tile/gf_list_tile.dart';


class UserView extends StatelessWidget {


  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Container(
        child: Column(
          children: [
            _buildHeader(),
          ],
        ),
      ),
    );
  }


  Widget _buildHeader(){
    return GFListTile(
        titleText:'Title',
        subTitleText:'Lorem ipsum dolor sit amet, consectetur adipiscing',
        icon: Icon(Icons.favorite)
    );
  }
}
