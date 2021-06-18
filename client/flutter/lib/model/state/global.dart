import 'package:get/get.dart';

class GlobalState extends GetxController{
  var selectedIndex = 0.obs;
  onItemTapped(int index) => selectedIndex.value = index;
}