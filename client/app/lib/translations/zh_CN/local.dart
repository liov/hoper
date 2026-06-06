import 'package:flutter/cupertino.dart';
import 'package:flutter/foundation.dart';

class MyLocalizationsDelegate
    extends LocalizationsDelegate<ZhCupertinoLocalizations> {
  const MyLocalizationsDelegate();

  @override
  bool isSupported(Locale locale) => locale.languageCode == 'zh' && locale.countryCode == 'CN';

  @override
  Future<ZhCupertinoLocalizations> load(Locale locale) =>
      ZhCupertinoLocalizations.load(locale);

  @override
  bool shouldReload(MyLocalizationsDelegate old) => false;

  @override
  String toString() => 'DefaultCupertinoLocalizations.delegate(zh-CN)';
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

  static Future<ZhCupertinoLocalizations> load(Locale locale) {
    return SynchronousFuture<ZhCupertinoLocalizations>(
        const ZhCupertinoLocalizations());
  }

  /// A [LocalizationsDelegate] that uses [DefaultCupertinoLocalizations.load]
  /// to create an instance of this class.
  static const LocalizationsDelegate<ZhCupertinoLocalizations> delegate =
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
  String get clearButtonLabel => '清除';

  @override
  String datePickerStandaloneMonth(int monthIndex) {
    return _months[monthIndex - 1];
  }

  @override
  String get lookUpButtonLabel => '查询';

  @override
  String get menuDismissLabel => '关闭';

  @override
  String get searchWebButtonLabel => '搜索';

  @override
  String get shareButtonLabel => '分享';

  @override
  String get backButtonLabel => '返回';

  @override
  String get cancelButtonLabel => '取消';

  @override
  String get collapsedHint => '折叠';

  @override
  String get expandedHint => '展开';

  @override
  String get expansionTileCollapsedHint => '折叠';

  @override
  String get expansionTileCollapsedTapHint => '折叠';

  @override
  String get expansionTileExpandedHint => '展开';

  @override
  String get expansionTileExpandedTapHint => '展开';
}
