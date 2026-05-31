import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/network/providers.dart';

final hasInstallationProvider = FutureProvider<bool>((ref) async {
  final client = ref.watch(apiClientProvider);
  try{
    final result = await client.get<Map<String, dynamic>>('api/github/installation');
    return result['installed'] as bool? ?? false;
  }catch (_) {
    return false;
  }
});