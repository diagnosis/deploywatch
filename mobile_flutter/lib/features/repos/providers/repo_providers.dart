import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/network/providers.dart';
import '../data/repos_service.dart';
import '../models/repo_model.dart';

final reposServiceProvider = Provider<ReposService>((ref) {
  return ReposService(ref.watch(apiClientProvider));
});

// Watched repos — used by both Feed and Repos screens
class WatchedReposNotifier extends AsyncNotifier<List<WatchedRepo>> {
  @override
  Future<List<WatchedRepo>> build() {
    return ref.read(reposServiceProvider).getWatchedRepos();
  }

  Future<void> watch(int repoId, String repoFullName, int installationId) async {
    await ref.read(reposServiceProvider).watchRepo(repoId, repoFullName, installationId);
    ref.invalidateSelf(); // refetch — like queryClient.invalidateQueries
  }

  Future<void> unwatch(int repoId) async {
    await ref.read(reposServiceProvider).unwatchRepo(repoId);
    ref.invalidateSelf();
  }
}

final watchedReposProvider =
AsyncNotifierProvider<WatchedReposNotifier, List<WatchedRepo>>(
  WatchedReposNotifier.new,
);

// GitHub repos provider
class GitHubReposNotifier extends AsyncNotifier<List<GitHubRepo>> {
  @override
  Future<List<GitHubRepo>> build() {
    return ref.read(reposServiceProvider).getGitHubRepos();
  }
}

final githubReposProvider =
AsyncNotifierProvider<GitHubReposNotifier, List<GitHubRepo>>(
  GitHubReposNotifier.new,
);