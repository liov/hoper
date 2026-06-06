import 'dart:convert';

class RbConnectionProfile {
  RbConnectionProfile({
    required this.id,
    required this.name,
    required this.signalUrl,
    this.direct = '',
    this.room = 'demo',
    this.path = '',
    this.agentPlatform = '',
  });

  final String id;
  final String name;
  final String signalUrl;
  final String direct;
  final String room;
  final String path;
  final String agentPlatform;

  factory RbConnectionProfile.fromJson(Map<String, dynamic> json) {
    return RbConnectionProfile(
      id: json['id'] as String? ?? '',
      name: json['name'] as String? ?? '',
      signalUrl: json['signalUrl'] as String? ?? '',
      direct: json['direct'] as String? ?? '',
      room: json['room'] as String? ?? 'demo',
      path: json['path'] as String? ?? '',
      agentPlatform: json['agentPlatform'] as String? ?? '',
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'name': name,
        'signalUrl': signalUrl,
        'direct': direct,
        'room': room,
        'path': path,
        if (agentPlatform.isNotEmpty) 'agentPlatform': agentPlatform,
      };

  static List<RbConnectionProfile> listFromJson(String raw) {
    if (raw.isEmpty) {
      return [];
    }
    final decoded = jsonDecode(raw);
    if (decoded is! List) {
      return [];
    }
    return decoded.map((e) => RbConnectionProfile.fromJson(Map<String, dynamic>.from(e as Map))).toList();
  }

  static String listToJson(List<RbConnectionProfile> list) => jsonEncode(list.map((e) => e.toJson()).toList());
}
