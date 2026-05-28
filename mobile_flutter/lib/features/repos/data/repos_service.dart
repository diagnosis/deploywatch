import '../../../core/network/api_client.dart';
import '../models/repo_model.dart';

class ReposService {
  final ApiClient _client;
  ReposService(this._client);

  Future<List<WatchedRepo>> getWatchedRepos() async {
    final data = await _client.get<Map<String, dynamic>>('/api/repos');
    final raw  = data['repos'] as List<dynamic>? ?? [];
    return raw.map((e) => WatchedRepo.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<void> watchRepo(int repoId, String repoFullName, int installationId) async {
    await _client.post<dynamic>('/api/repos/watch', {
      'repo_id':         repoId,
      'repo_full_name':  repoFullName,
      'installation_id': installationId,
    });
  }

  Future<void> unwatchRepo(int repoId) async {
    await _client.delete<dynamic>('/api/repos/watch/$repoId');
  }
  Future<List<GitHubRepo>> getGitHubRepos() async {
    final data = await _client.get<List<dynamic>>('/api/github/repos');
    return data.map((e) => GitHubRepo.fromJson(e as Map<String, dynamic>)).toList();
  }
}