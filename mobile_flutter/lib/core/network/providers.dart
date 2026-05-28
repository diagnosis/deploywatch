import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'token_store.dart';
import 'api_client.dart';

final tokenStoreProvider = Provider<TokenStore>((ref) {
  return TokenStore();
});

final apiClientProvider = Provider<ApiClient>((ref) {
  final tokenStore = ref.watch(tokenStoreProvider);
  return ApiClient(tokenStore);
});