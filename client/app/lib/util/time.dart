
DateTime getDateTime(int seconds,int nanos) {
  return DateTime.fromMillisecondsSinceEpoch(seconds * 1000 + (nanos/1000000) as int);
}