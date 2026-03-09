import 'package:flutter/cupertino.dart';
import 'package:flutter/foundation.dart';

class MyLocalizationsDelegate
    extends LocalizationsDelegate<CupertinoLocalizations> {
  const MyLocalizationsDelegate();

  @override
  bool isSupported(Locale locale) => locale.languageCode == 'zh';

  @override
  Future<CupertinoLocalizations> load(Locale locale) =>
      ZhCupertinoLocalizations.load(locale);

  @override
  bool shouldReload(MyLocalizationsDelegate old) => false;

  @override
  String toString() => 'DefaultCupertinoLocalizations.delegate(zh)';
}


class ZhCupertinoLocalizations implements CupertinoLocalizations {
  const ZhCupertinoLocalizations();

  static const List<String> _shortWeekdays = <String>[
    '星期一',
    '星期二',
    '星期三',
    '星期四',
    '星期五',
    '星期六',
    '星期日',
  ];

  static const List<String> _months = <String>[
    '1月',
    '2月',
    '3月',
    '4月',
    '5月',
    '6月',
    '7月',
    '8月',
    '9月',
    '10月',
    '11月',
    '12月',
  ];


  @override
  String datePickerYear(int yearIndex) => yearIndex.toString();

  @override
  String datePickerMonth(int monthIndex) => _months[monthIndex - 1];

  @override
  String datePickerDayOfMonth(int dayIndex, [int? weekDay]) => dayIndex.toString();

  @override
  String datePickerHour(int hour) => hour.toString();

  @override
  String datePickerHourSemanticsLabel(int hour) => "$hour o'clock";

  @override
  String datePickerMinute(int minute) => minute.toString().padLeft(2, '0');

  @override
  String datePickerMinuteSemanticsLabel(int minute) {
    if (minute == 1) return '1 分';
    return '$minute 分';
  }

  @override
  String datePickerMediumDate(DateTime date) {
    return '${_shortWeekdays[date.weekday - DateTime.monday]} '
        '${_months[date.month - DateTime.january]} '
        '${date.day.toString().padRight(2)}';
  }

  @override
  DatePickerDateOrder get datePickerDateOrder => DatePickerDateOrder.mdy;

  @override
  DatePickerDateTimeOrder get datePickerDateTimeOrder =>
      DatePickerDateTimeOrder.date_time_dayPeriod;

  @override
  String get anteMeridiemAbbreviation => '上午';

  @override
  String get postMeridiemAbbreviation => '下午';

  @override
  String get todayLabel => '今天';

  @override
  String get alertDialogLabel => 'Alert';

  @override
  String timerPickerHour(int hour) => hour.toString();

  @override
  String timerPickerMinute(int minute) => minute.toString();

  @override
  String timerPickerSecond(int second) => second.toString();

  @override
  String timerPickerHourLabel(int hour) =>  '时';

  @override
  String timerPickerMinuteLabel(int minute) => '分';

  @override
  String timerPickerSecondLabel(int second) => '秒';

  @override
  String get cutButtonLabel => '剪贴';

  @override
  String get copyButtonLabel => '复制';

  @override
  String get pasteButtonLabel => '粘贴';

  @override
  String get selectAllButtonLabel => '选择全部';

  static Future<CupertinoLocalizations> load(Locale locale) {
    return SynchronousFuture<CupertinoLocalizations>(
        const ZhCupertinoLocalizations());
  }

  /// A [LocalizationsDelegate] that uses [DefaultCupertinoLocalizations.load]
  /// to create an instance of this class.
  static const LocalizationsDelegate<CupertinoLocalizations> delegate =
  MyLocalizationsDelegate();

  @override
  String get modalBarrierDismissLabel => '取消模态窗';

  @override
  String get searchTextFieldPlaceholderLabel => '搜索占位符';

  @override
  String tabSemanticsLabel({required int tabIndex, required int tabCount}) {
    return '标签$tabIndex/$tabCount';
  }

  @override
  List<String> get timerPickerHourLabels => List.generate(24, (index) => (index+1).toString(),growable: false);

  @override
  List<String> get timerPickerMinuteLabels => List.generate(60, (index) => index.toString(),growable: false);

  @override
  List<String> get timerPickerSecondLabels => List.generate(60, (index) => index.toString(),growable: false);

  @override
  String get noSpellCheckReplacementsLabel => '';

  @override
  // TODO: implement clearButtonLabel
  String get clearButtonLabel => '清除';

  @override
  String datePickerStandaloneMonth(int monthIndex) {
    // TODO: implement datePickerStandaloneMonth
    return _months[monthIndex - 1];
  }

  @override
  // TODO: implement lookUpButtonLabel
  String get lookUpButtonLabel => '查询';

  @override
  // TODO: implement menuDismissLabel
  String get menuDismissLabel => '关闭';

  @override
  // TODO: implement searchWebButtonLabel
  String get searchWebButtonLabel => '搜索';

  @override
  // TODO: implement shareButtonLabel
  String get shareButtonLabel => '分享';

  @override
  // TODO: implement backButtonLabel
  String get backButtonLabel => '返回';

  @override
  // TODO: implement cancelButtonLabel
  String get cancelButtonLabel => '取消';

  @override
  // TODO: implement collapsedHint
  String get collapsedHint => '折叠';

  @override
  // TODO: implement expandedHint
  String get expandedHint => '展开';

  @override
  // TODO: implement expansionTileCollapsedHint
  String get expansionTileCollapsedHint => '折叠';

  @override
  // TODO: implement expansionTileCollapsedTapHint
  String get expansionTileCollapsedTapHint => '折叠';

  @override
  // TODO: implement expansionTileExpandedHint
  String get expansionTileExpandedHint => '展开';

  @override
  // TODO: implement expansionTileExpandedTapHint
  String get expansionTileExpandedTapHint => '展开';
}
