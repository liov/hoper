
dynamic jsonObject({Map<String, dynamic>? json}) => JsonObject._(json ?? {});

class JsonObject {
  final Map<String, dynamic> map;

  const JsonObject._(this.map);

  @override
  dynamic noSuchMethod(Invocation i) {
    if (i.isGetter) {
      return map[i.memberName.name];
    } else if (i.isSetter) {
      map[i.memberName.name] = i.positionalArguments.first;
    } else {
      return super.noSuchMethod(i);
    }
  }

  @override
  String toString() => map.toString();
}

extension _SymbolExt on Symbol {
  String get name {
    final s = toString();
    final symbol = s.substring(8, s.length - 2);
    return symbol.endsWith('=') ? symbol.substring(0, symbol.length - 1) : symbol;
  }
}