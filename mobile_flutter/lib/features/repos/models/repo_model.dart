class WatchedRepo {
  final String id;
  final int repoId;
  final String repoFullName;
  final String createdAt;

  const WatchedRepo({
    required this.id,
    required this.repoId,
    required this.repoFullName,
    required this.createdAt,
  });

  factory WatchedRepo.fromJson(Map<String, dynamic> json) {
    return WatchedRepo(
      id:           json['id']?.toString() ?? '',
      repoId:       json['repo_id'] as int,
      repoFullName: json['repo_full_name'] as String,
      createdAt:    json['created_at']?.toString() ?? '',
    );
  }
}

class GitHubRepo {
  final int id;
  final String name;
  final String fullName;
  final int installationId;

  const GitHubRepo({
    required this.id,
    required this.name,
    required this.fullName,
    required this.installationId,
  });

  factory GitHubRepo.fromJson(Map<String, dynamic> json) {
    return GitHubRepo(
      id:             json['id']              as int,
      name:           json['name']            as String,
      fullName:       json['full_name']       as String,
      installationId: json['installation_id'] as int,
    );
  }
}