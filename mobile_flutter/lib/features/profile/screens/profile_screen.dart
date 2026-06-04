import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:url_launcher/url_launcher.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/network/providers.dart';
import '../../auth/providers/auth_providers.dart';

// User model
class User {
  final String id;
  final String login;
  final String? name;
  final String? email;

  const User({
    required this.id,
    required this.login,
    this.name,
    this.email,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id:    json['id']    as String,
      login: json['login'] as String,
      name:  json['name']  as String?,
      email: json['email'] as String?,
    );
  }
}

// User provider
final meProvider = FutureProvider<User>((ref) async {
  final client = ref.watch(apiClientProvider);
  final data   = await client.get<Map<String, dynamic>>('/api/auth/me');
  return User.fromJson(data['user'] as Map<String, dynamic>);
});
final deleteAccountProvider = FutureProvider.autoDispose<void>((ref) async {
  final client = ref.watch(apiClientProvider);
  await client.delete<Map<String, dynamic>>('/api/auth/account');
});

class ProfileScreen extends ConsumerWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final userAsync = ref.watch(meProvider);

    return Scaffold(
      backgroundColor: AppColors.background,
      body: SafeArea(
        child: SingleChildScrollView(
          child: Column(
            children: [
              // Header
              Container(
                width: double.infinity,
                padding: const EdgeInsets.symmetric(
                  horizontal: 16,
                  vertical: 12,
                ),
                decoration: const BoxDecoration(
                  border: Border(bottom: BorderSide(color: AppColors.border)),
                ),
                child: const Text(
                  'Profile',
                  style: TextStyle(
                    color: AppColors.textPrimary,
                    fontSize: 20,
                    fontWeight: FontWeight.w700,
                  ),
                ),
              ),

              // Avatar
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 32),
                child: Column(
                  children: [
                    Container(
                      width: 80,
                      height: 80,
                      decoration: BoxDecoration(
                        color: AppColors.primaryDim,
                        shape: BoxShape.circle,
                        border: Border.all(color: AppColors.borderAlt, width: 2),
                      ),
                      child: const Icon(
                        Icons.person,
                        size: 36,
                        color: AppColors.primary,
                      ),
                    ),
                    const SizedBox(height: 12),
                    userAsync.when(
                      loading: () => const CircularProgressIndicator(
                        color: AppColors.primary,
                      ),
                      error: (_, __) => const Text(
                        '—',
                        style: TextStyle(color: AppColors.textPrimary),
                      ),
                      data: (user) => Column(
                        children: [
                          Text(
                            user.login,
                            style: const TextStyle(
                              color: AppColors.textPrimary,
                              fontSize: 18,
                              fontWeight: FontWeight.w700,
                            ),
                          ),
                          if (user.name != null) ...[
                            const SizedBox(height: 4),
                            Text(
                              user.name!,
                              style: const TextStyle(
                                color: AppColors.textSecondary,
                                fontSize: 14,
                              ),
                            ),
                          ],
                          if (user.email != null) ...[
                            const SizedBox(height: 2),
                            Text(
                              user.email!,
                              style: const TextStyle(
                                color: AppColors.textDim,
                                fontSize: 13,
                              ),
                            ),
                          ],
                        ],
                      ),
                    ),
                  ],
                ),
              ),

              // Info card
              Container(
                margin: const EdgeInsets.symmetric(horizontal: 16),
                decoration: BoxDecoration(
                  color: AppColors.surface,
                  borderRadius: BorderRadius.circular(12),
                  border: Border.all(color: AppColors.border),
                ),
                child: userAsync.maybeWhen(
                  data: (user) => Column(
                    children: [
                      _InfoRow(
                        icon: Icons.code,
                        label: 'GitHub',
                        value: '@${user.login}',
                        hasBorder: user.email != null,
                      ),
                      if (user.email != null)
                        _InfoRow(
                          icon: Icons.mail_outline,
                          label: 'Email',
                          value: user.email!,
                          hasBorder: false,
                        ),
                    ],
                  ),
                  orElse: () => const SizedBox.shrink(),
                ),
              ),

              const SizedBox(height: 24),

              // Sign out button
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16),
                child: SizedBox(
                  width: double.infinity,
                  child: ElevatedButton.icon(
                    onPressed: () async {
                      await ref.read(authProvider.notifier).logout();
                      if (context.mounted) context.go('/login');
                    },
                    icon: const Icon(Icons.logout, size: 18),
                    label: const Text(
                      'Sign out',
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: AppColors.error,
                      foregroundColor: Colors.white,
                      padding: const EdgeInsets.symmetric(vertical: 16),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                  ),
                ),
              ),

              const SizedBox(height: 12),


// Delete account button
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16),
                child: SizedBox(
                  width: double.infinity,
                  child: OutlinedButton.icon(
                    onPressed: () async {
                      // Show confirmation dialog
                      final confirm = await showDialog<bool>(
                        context: context,
                        builder: (ctx) => AlertDialog(
                          backgroundColor: AppColors.surface,
                          title: const Text('Delete Account',
                              style: TextStyle(color: AppColors.textPrimary)),
                          content: const Text(
                              'This will permanently delete your account and all data. This cannot be undone.',
                              style: TextStyle(color: AppColors.textSecondary)),
                          actions: [
                            TextButton(
                              onPressed: () => Navigator.pop(ctx, false),
                              child: const Text('Cancel',
                                  style: TextStyle(color: AppColors.textSecondary)),
                            ),
                            TextButton(
                              onPressed: () => Navigator.pop(ctx, true),
                              child: const Text('Delete',
                                  style: TextStyle(color: AppColors.error)),
                            ),
                          ],
                        ),
                      );
                      if (confirm != true) return;
                      final client = ref.read(apiClientProvider);
                      await client.delete<Map<String, dynamic>>('/api/auth/account');
                      await ref.read(authProvider.notifier).logout();
                      if (context.mounted) context.go('/login');
                    },
                    icon: const Icon(Icons.delete_outline, size: 18, color: AppColors.error),
                    label: const Text('Delete Account',
                        style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600, color: AppColors.error)),
                    style: OutlinedButton.styleFrom(
                      side: const BorderSide(color: AppColors.error),
                      padding: const EdgeInsets.symmetric(vertical: 16),
                      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                    ),
                  ),
                ),
              ),
              const SizedBox(height: 24),
              const Text(
                'deploywatch v1.0.0',
                style: TextStyle(
                  color: AppColors.textDim,
                  fontSize: 12,
                ),
              ),
              GestureDetector(
                onTap: () async {
                  final uri = Uri.parse('https://gist.githubusercontent.com/diagnosis/7904a39226032ee521d7f93b026f58d0/raw/6e01f56eb2602b480d2740b6070fc553db648a15/privacy-policy.md');
                  await launchUrl(uri, mode: LaunchMode.externalApplication);
                },
                child: const Text(
                  'Privacy Policy',
                  style: TextStyle(
                    color: AppColors.primary,
                    fontSize: 12,
                    decoration: TextDecoration.underline,
                  ),
                ),
              ),



              const SizedBox(height: 24),
            ],
          ),
        ),
      ),
    );
  }
}

class _InfoRow extends StatelessWidget {
  final IconData icon;
  final String label;
  final String value;
  final bool hasBorder;

  const _InfoRow({
    required this.icon,
    required this.label,
    required this.value,
    required this.hasBorder,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        border: hasBorder
            ? const Border(bottom: BorderSide(color: AppColors.border))
            : null,
      ),
      child: Row(
        children: [
          Icon(icon, size: 18, color: AppColors.textSecondary),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              label,
              style: const TextStyle(
                color: AppColors.textSecondary,
                fontSize: 14,
              ),
            ),
          ),
          Text(
            value,
            style: const TextStyle(
              color: AppColors.textPrimary,
              fontSize: 14,
            ),
          ),
        ],
      ),
    );
  }
}