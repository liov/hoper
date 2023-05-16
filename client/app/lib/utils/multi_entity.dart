mixin MultiEntity<T> {
  final Map<String, T> entityMap = Map();

  newEntity(T t, String tag) async {
    entityMap[tag] = t;
  }

  T? getEntity(String tag) {
    return entityMap[tag];
  }
}
