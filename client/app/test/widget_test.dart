import 'package:app/app.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  testWidgets('AppRoot smoke test', (WidgetTester tester) async {
    await tester.pumpWidget(const ProviderScope(child: AppRoot()));
    await tester.pump();
    expect(find.byType(AppRoot), findsOneWidget);
  });
}
