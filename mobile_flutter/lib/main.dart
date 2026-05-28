import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'core/theme/app_theme.dart';
import 'core/router/router.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  runApp(
    const ProviderScope(
      child: DeployWatchApp(),
    ),
  );
}

// ConsumerWidget = a widget that can read Riverpod providers
// Use this instead of StatelessWidget when you need providers
class DeployWatchApp extends ConsumerWidget {
  const DeployWatchApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // ref.watch() here means: rebuild this widget when router changes
    // router changes when auth state changes → redirects automatically
    final router = ref.watch(routerProvider);

    return MaterialApp.router(
      title: 'DeployWatch',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.dark,
      routerConfig: router,
    );
  }
}