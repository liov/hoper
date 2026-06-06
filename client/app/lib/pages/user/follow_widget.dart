import 'package:app/providers/providers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

class FollowWidget extends ConsumerWidget {
  const FollowWidget({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    ref.watch(followProvider);
    return Container();
  }
}
