import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/network/providers.dart';
import '../../feed/providers/feed_providers.dart';
import '../../github/github_providers.dart';
import '../../profile/screens/profile_screen.dart';
import '../../repos/providers/repo_providers.dart';
import '../data/auth_service.dart';

// AuthService provider — depends on tokenStore
final authServiceProvider = Provider<AuthService>((ref) {
  final tokenStore = ref.watch(tokenStoreProvider);
  return AuthService(tokenStore);
});

// This tells the whole app whether the user is logged in.
// It's an AsyncNotifier — like useQuery but with methods.
// 'state' holds an AsyncValue<bool>:
//   AsyncValue.loading()  → still checking
//   AsyncValue.data(true) → logged in
//   AsyncValue.data(false)→ not logged in
class AuthNotifier extends AsyncNotifier<bool> {
  @override
  Future<bool> build() async {
    final tokenStore = ref.read(tokenStoreProvider);
    final hasToken = await tokenStore.hasToken();
    if (hasToken) {
      Future.microtask(() =>
          ref.read(notificationServiceProvider).initialize()
      );
    }
    return hasToken;
  }

  Future<bool> login() async {
    final authService = ref.read(authServiceProvider);
    final success = await authService.loginWithGithub();
    if (success) {
      ref.invalidate(watchedReposProvider);
      ref.invalidate(eventsProvider);
      ref.invalidate(hasInstallationProvider);
      ref.invalidate(meProvider);
      state = const AsyncValue.data(true);
      await ref.read(notificationServiceProvider).initialize();
    }
    return success;
  }

  Future<void> logout() async {
    final authService = ref.read(authServiceProvider);
    await authService.logout();
    state = const AsyncValue.data(false);
  }
}

final authProvider = AsyncNotifierProvider<AuthNotifier, bool>(
  AuthNotifier.new,
);