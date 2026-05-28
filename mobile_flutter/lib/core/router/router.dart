import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../features/auth/providers/auth_providers.dart';
import '../../features/auth/screens/login_screen.dart';
import '../../features/feed/screens/feed_screen.dart';
import '../../features/repos/screens/repos_screen.dart';
import '../../features/profile/screens/profile_screen.dart';




// Shell screen — holds the bottom tab bar
// All tab screens live inside this
class ShellScreen extends StatelessWidget {
  final Widget child;
  const ShellScreen({super.key, required this.child});

  int _locationToIndex(String location) {
    if (location.startsWith('/repos'))   return 1;
    if (location.startsWith('/profile')) return 2;
    return 0; // feed
  }

  @override
  Widget build(BuildContext context) {
    final location = GoRouterState.of(context).uri.toString();
    return Scaffold(
      backgroundColor: const Color(0xFF080D12),
      body: child,
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: _locationToIndex(location),
        onTap: (index) {
          switch (index) {
            case 0: context.go('/feed');    break;
            case 1: context.go('/repos');   break;
            case 2: context.go('/profile'); break;
          }
        },
        items: const [
          BottomNavigationBarItem(
            icon: Icon(Icons.flash_on_outlined),
            activeIcon: Icon(Icons.flash_on),
            label: 'Feed',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.account_tree_outlined),
            activeIcon: Icon(Icons.account_tree),
            label: 'Repos',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.person_outline),
            activeIcon: Icon(Icons.person),
            label: 'Profile',
          ),
        ],
      ),
    );
  }
}

// The router provider
final routerProvider = Provider<GoRouter>((ref) {
  // Listen to auth state — router rebuilds when login state changes
  final authState = ref.watch(authProvider);

  return GoRouter(
    initialLocation: '/feed',
    // redirect runs before every navigation
    // This is the auth guard
    redirect: (context, state) {
      final isLoggedIn = authState.valueOrNull ?? false;
      final isLoggingIn = state.matchedLocation == '/login';

      if (!isLoggedIn && !isLoggingIn) return '/login';
      if (isLoggedIn && isLoggingIn)   return '/feed';
      return null; // no redirect
    },
    routes: [
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginScreen(),
      ),
      ShellRoute(
        builder: (context, state, child) => ShellScreen(child: child),
        routes: [
          GoRoute(path: '/feed',    builder: (_, __) => const FeedScreen()),
          GoRoute(path: '/repos',   builder: (_, __) => const ReposScreen()),
          GoRoute(path: '/profile', builder: (_, __) => const ProfileScreen()),
        ],
      ),
    ],
  );
});