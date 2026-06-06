import 'package:app/remotebrowse/rb_connection_profile.dart';
import 'package:shared_preferences/shared_preferences.dart';

class RbConnectionStore {
  static const _listKey = 'rb_connection_profiles_v1';
  static const _lastIdKey = 'rb_connection_last_id_v1';

  static Future<List<RbConnectionProfile>> loadAll() async {
    final prefs = await SharedPreferences.getInstance();
    return RbConnectionProfile.listFromJson(prefs.getString(_listKey) ?? '');
  }

  static Future<String?> loadLastId() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_lastIdKey);
  }

  static Future<void> persist(List<RbConnectionProfile> list, {String? lastId}) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_listKey, RbConnectionProfile.listToJson(list));
    if (lastId != null) {
      await prefs.setString(_lastIdKey, lastId);
    }
  }

  static Future<void> remove(String id) async {
    final list = await loadAll()..removeWhere((p) => p.id == id);
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_listKey, RbConnectionProfile.listToJson(list));
    final last = await loadLastId();
    if (last != id) {
      return;
    }
    if (list.isEmpty) {
      await prefs.remove(_lastIdKey);
    } else {
      await prefs.setString(_lastIdKey, list.first.id);
    }
  }
}
