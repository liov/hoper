import 'dart:io';

String rbPlatformName() => Platform.operatingSystem;

bool rbHasIpv6() => Platform.isAndroid || Platform.isIOS;
