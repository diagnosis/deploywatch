import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_colors.dart';
import '../providers/repo_providers.dart';
import '../models/repo_model.dart';

class ReposScreen extends ConsumerWidget {
  const ReposScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final reposAsync = ref.watch(watchedReposProvider);

    return Scaffold(
      backgroundColor: AppColors.background,
      body: SafeArea(
        child: Column(
          children: [
            // Header
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
              decoration: const BoxDecoration(
                border: Border(bottom: BorderSide(color: AppColors.border)),
              ),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  const Text(
                    'Watched Repos',
                    style: TextStyle(
                      color: AppColors.textPrimary,
                      fontSize: 20,
                      fontWeight: FontWeight.w700,
                    ),
                  ),
                  // Add button
                  GestureDetector(
                    onTap: () => _showAddRepoSheet(context, ref),
                    child: Container(
                      width: 32,
                      height: 32,
                      decoration: BoxDecoration(
                        color: AppColors.primaryDim,
                        borderRadius: BorderRadius.circular(8),
                        border: Border.all(color: AppColors.borderAlt),
                      ),
                      child: const Icon(
                        Icons.add,
                        size: 18,
                        color: AppColors.primary,
                      ),
                    ),
                  ),
                ],
              ),
            ),

            // Repos list
            Expanded(
              child: reposAsync.when(
                loading: () => const Center(
                  child: CircularProgressIndicator(color: AppColors.primary),
                ),
                error: (err, _) => Center(
                  child: Text(
                    'Error: $err',
                    style: const TextStyle(color: AppColors.error),
                  ),
                ),
                data: (repos) {
                  if (repos.isEmpty) {
                    return Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          const Icon(
                            Icons.account_tree_outlined,
                            size: 48,
                            color: AppColors.textDim,
                          ),
                          const SizedBox(height: 16),
                          const Text(
                            'No repos watched',
                            style: TextStyle(color: AppColors.textSecondary),
                          ),
                          const SizedBox(height: 12),
                          GestureDetector(
                            onTap: () => _showAddRepoSheet(context, ref),
                            child: const Text(
                              '+ Watch a repo',
                              style: TextStyle(color: AppColors.primary),
                            ),
                          ),
                        ],
                      ),
                    );
                  }

                  return RefreshIndicator(
                    color: AppColors.primary,
                    onRefresh: () => ref.refresh(watchedReposProvider.future),
                    child: ListView.builder(
                      itemCount: repos.length,
                      itemBuilder: (context, index) {
                        final repo = repos[index];
                        return _RepoItem(
                          repo: repo,
                          onUnwatch: () => _confirmUnwatch(context, ref, repo),
                        );
                      },
                    ),
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _confirmUnwatch(BuildContext context, WidgetRef ref, WatchedRepo repo) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: AppColors.surface,
        title: const Text(
          'Unwatch',
          style: TextStyle(color: AppColors.textPrimary),
        ),
        content: Text(
          'Stop watching ${repo.repoFullName}?',
          style: const TextStyle(color: AppColors.textSecondary),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text(
              'Cancel',
              style: TextStyle(color: AppColors.textSecondary),
            ),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(ctx);
              ref.read(watchedReposProvider.notifier).unwatch(repo.repoId);
            },
            child: const Text(
              'Unwatch',
              style: TextStyle(color: AppColors.error),
            ),
          ),
        ],
      ),
    );
  }

  void _showAddRepoSheet(BuildContext context, WidgetRef ref) {
    showModalBottomSheet(
      context: context,
      backgroundColor: AppColors.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (ctx) => _AddRepoSheet(),
    );
  }
}

class _RepoItem extends StatelessWidget {
  final WatchedRepo repo;
  final VoidCallback onUnwatch;

  const _RepoItem({required this.repo, required this.onUnwatch});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
      decoration: const BoxDecoration(
        border: Border(bottom: BorderSide(color: AppColors.border)),
      ),
      child: Row(
        children: [
          const Icon(
            Icons.account_tree_outlined,
            size: 16,
            color: AppColors.primary,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              repo.repoFullName,
              style: const TextStyle(
                color: AppColors.textPrimary,
                fontSize: 14,
                fontFamily: 'monospace',
              ),
            ),
          ),
          GestureDetector(
            onTap: onUnwatch,
            child: const Icon(
              Icons.delete_outline,
              size: 16,
              color: AppColors.error,
            ),
          ),
        ],
      ),
    );
  }
}

// Bottom sheet for adding repos
class _AddRepoSheet extends ConsumerWidget {
  const _AddRepoSheet();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final githubReposAsync = ref.watch(githubReposProvider);
    final watchedRepos     = ref.watch(watchedReposProvider).valueOrNull ?? [];

    return Padding(
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text(
                'Watch a Repo',
                style: TextStyle(
                  color: AppColors.textPrimary,
                  fontSize: 18,
                  fontWeight: FontWeight.w700,
                ),
              ),
              GestureDetector(
                onTap: () => Navigator.pop(context),
                child: const Icon(
                  Icons.close,
                  color: AppColors.textSecondary,
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          Expanded(
            child: githubReposAsync.when(
              loading: () => const Center(
                child: CircularProgressIndicator(color: AppColors.primary),
              ),
              error: (err, _) => Center(
                child: Text(
                  'Error: $err',
                  style: const TextStyle(color: AppColors.error),
                ),
              ),
              data: (repos) => ListView.builder(
                itemCount: repos.length,
                itemBuilder: (context, index) {
                  final repo = repos[index];
                  final alreadyWatched =
                  watchedRepos.any((r) => r.repoId == repo.id);
                  return Opacity(
                    opacity: alreadyWatched ? 0.4 : 1,
                    child: GestureDetector(
                      onTap: alreadyWatched
                          ? null
                          : () {
                        ref
                            .read(watchedReposProvider.notifier)
                            .watch(
                          repo.id,
                          repo.fullName,
                          repo.installationId,
                        );
                        Navigator.pop(context);
                      },
                      child: Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 16,
                          vertical: 12,
                        ),
                        decoration: const BoxDecoration(
                          border: Border(
                            bottom: BorderSide(color: AppColors.border),
                          ),
                        ),
                        child: Row(
                          children: [
                            const Icon(
                              Icons.account_tree_outlined,
                              size: 16,
                              color: AppColors.primary,
                            ),
                            const SizedBox(width: 12),
                            Expanded(
                              child: Text(
                                repo.fullName,
                                style: const TextStyle(
                                  color: AppColors.textPrimary,
                                  fontSize: 14,
                                  fontFamily: 'monospace',
                                ),
                              ),
                            ),
                            if (alreadyWatched)
                              const Text(
                                'watching',
                                style: TextStyle(
                                  color: AppColors.textDim,
                                  fontSize: 12,
                                ),
                              ),
                          ],
                        ),
                      ),
                    ),
                  );
                },
              ),
            ),
          ),
        ],
      ),
    );
  }
}