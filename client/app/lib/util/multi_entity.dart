mixin MultiEntity<T> {
  final Map<String, T> entityMap = {};

  Future<void> newEntity(T t, String tag) async {
    entityMap[tag] = t;
  }

  T? getEntity(String tag) {
    return entityMap[tag];
  }
}
