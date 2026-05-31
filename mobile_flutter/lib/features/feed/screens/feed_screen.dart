import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile_flutter/features/github/github_providers.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/network/providers.dart';
import '../../repos/providers/repo_providers.dart';
import '../providers/feed_providers.dart';
import '../models/event_model.dart';
import 'dart:convert';
import 'package:url_launcher/url_launcher.dart';

// ── Helpers ────────────────────────────────────────────────────

String timeAgo(String dateStr) {
  final date    = DateTime.parse(dateStr);
  final seconds = DateTime.now().difference(date).inSeconds;
  if (seconds < 60)   return '${seconds}s ago';
  if (seconds < 3600) return '${seconds ~/ 60}m ago';  // ~/ = integer division
  if (seconds < 86400) return '${seconds ~/ 3600}h ago';
  return '${seconds ~/ 86400}d ago';
}

Color eventColor(String type) {
  switch (type) {
    case 'push':                return AppColors.pushEvent;
    case 'pull_request':        return AppColors.prEvent;
    case 'pull_request_review': return AppColors.reviewEvent;
    default:                    return AppColors.textSecondary;
  }
}

IconData eventIcon(String type) {
  switch (type) {
    case 'push':         return Icons.commit;
    case 'pull_request': return Icons.merge_type;
    default:             return Icons.flash_on;
  }
}

// ── Event Card ─────────────────────────────────────────────────

void _openGitHub(Event event, String repoFullName) async {
  final url = _getGitHubUrl(event, repoFullName);
  final uri = Uri.parse(url);
  if (await canLaunchUrl(uri)) {
    await launchUrl(uri, mode: LaunchMode.externalApplication);
  }
}

String _getGitHubUrl(Event event, String repoFullName) {
  Map<String, dynamic> p = {};
  try {
    p = json.decode(event.payload) as Map<String, dynamic>;
  } catch (_) {}

  switch (event.eventType) {
    case 'push':
      return 'https://github.com/$repoFullName/commit/${p['after'] ?? ''}';
    case 'pull_request':
      return 'https://github.com/$repoFullName/pull/${p['number'] ?? ''}';
    case 'pull_request_review':
      final prNumber = (p['pull_request'] as Map?)?.cast<String,dynamic>()['number'];
      return 'https://github.com/$repoFullName/pull/$prNumber';
    case 'create':
      return 'https://github.com/$repoFullName/tree/${p['ref'] ?? ''}';
    default:
      return 'https://github.com/$repoFullName';
  }
}

class EventCard extends StatelessWidget {
  final Event event;
  final String repoFullName;

  const EventCard({
    super.key,
    required this.event,
    required this.repoFullName,
  });


  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
      decoration: const BoxDecoration(
        border: Border(
          bottom: BorderSide(color: AppColors.border),
        ),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Event type icon
          Container(
            width: 32,
            height: 32,
            decoration: BoxDecoration(
              color: AppColors.surface,
              borderRadius: BorderRadius.circular(8),
              border: Border.all(color: AppColors.border),
            ),
            child: Icon(
              eventIcon(event.eventType),
              size: 14,
              color: eventColor(event.eventType),
            ),
          ),

          const SizedBox(width: 12),

          // Event details
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Event type + action badge
                Row(
                  children: [
                    Text(
                      event.eventType,
                      style: const TextStyle(
                        color: AppColors.textPrimary,
                        fontSize: 14,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    if (event.action != null) ...[
                      const SizedBox(width: 8),
                      Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 8,
                          vertical: 2,
                        ),
                        decoration: BoxDecoration(
                          color: AppColors.primaryDim,
                          borderRadius: BorderRadius.circular(10),
                          border: Border.all(color: AppColors.borderAlt),
                        ),
                        child: Text(
                          event.action!,
                          style: const TextStyle(
                            color: AppColors.primary,
                            fontSize: 11,
                          ),
                        ),
                      ),
                    ],
                  ],
                ),


                const SizedBox(height: 2),

                // Actor + time
                Row(
                  children: [
                    const Icon(
                      Icons.person_outline,
                      size: 11,
                      color: AppColors.textDim,
                    ),
                    const SizedBox(width: 4),
                    Text(
                      event.actorLogin,
                      style: const TextStyle(
                        color: AppColors.textSecondary,
                        fontSize: 12,
                      ),
                    ),
                    const SizedBox(width: 6),
                    const Text(
                      '·',
                      style: TextStyle(
                        color: AppColors.textDim,
                        fontSize: 12,
                      ),
                    ),
                    const SizedBox(width: 6),
                    Text(
                      timeAgo(event.receivedAt),
                      style: const TextStyle(
                        color: AppColors.textDim,
                        fontSize: 12,
                      ),
                    ),
                  ],
                ),

                const SizedBox(height: 4),

                // Repo name
                Text(
                  repoFullName,
                  style: const TextStyle(
                    color: AppColors.textDim,
                    fontSize: 11,
                    fontFamily: 'monospace',
                  ),
                ),
                GestureDetector(
                  onTap: () => _openGitHub(event, repoFullName),
                  child: Padding(
                    padding: const EdgeInsets.only(top: 6),
                    child: Row(
                      children: const [
                        Icon(Icons.open_in_new, size: 11, color: AppColors.primary),
                        SizedBox(width: 4),
                        Text(
                          'View on GitHub',
                          style: TextStyle(
                            color: AppColors.primary,
                            fontSize: 11,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ],
            ),
          ),
        ],

      ),

    );
  }
}

// ── Feed Screen ────────────────────────────────────────────────

// ConsumerStatefulWidget because we need local state (dropdown open/closed)
// AND we need to read providers
class FeedScreen extends ConsumerStatefulWidget {
  const FeedScreen({super.key});
  @override
  ConsumerState<FeedScreen> createState() => _FeedScreenState();
}

class _FeedScreenState extends ConsumerState<FeedScreen> {
  bool _showRepoPicker = false;
  Timer? _timer;


  @override
  void initState(){
    super.initState();
    _timer = Timer.periodic(const Duration(seconds: 30), (_){
      if (mounted) setState(() {});
    });
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    // Watch providers — widget rebuilds when these change
    final eventsAsync  = ref.watch(eventsProvider);
    final reposAsync   = ref.watch(watchedReposProvider);
    final selectedRepo = ref.watch(selectedRepoIdProvider);
    final currentPage  = ref.watch(currentPageProvider);
    final hasInstallation = ref.watch(hasInstallationProvider);
    ref.watch(sseNotifierProvider);
    final repos = reposAsync.valueOrNull ?? [];

    final selectedRepoName = selectedRepo == null
        ? 'All repos'
        : repos
        .where((r) => r.repoId == selectedRepo)
        .map((r) => r.repoFullName.split('/').last)
        .firstOrNull ?? 'All repos';

    return Scaffold(
      backgroundColor: AppColors.background,
      body: SafeArea(
        child: Column(
          children: [
            // ── Header ──────────────────────────────────────
            Container(
              padding: const EdgeInsets.symmetric(
                horizontal: 16,
                vertical: 12,
              ),
              decoration: const BoxDecoration(
                border: Border(bottom: BorderSide(color: AppColors.border)),
              ),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Row(
                    children: [
                      const Icon(
                        Icons.flash_on,
                        color: AppColors.primary,
                        size: 18,
                      ),
                      const SizedBox(width: 8),
                      const Text(
                        'deploywatch',
                        style: TextStyle(
                          color: AppColors.textPrimary,
                          fontSize: 16,
                          fontWeight: FontWeight.w700,
                        ),
                      ),
                    ],
                  ),
                  // Live badge
                  Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 4,
                    ),
                    decoration: BoxDecoration(
                      color: AppColors.primaryDim,
                      borderRadius: BorderRadius.circular(10),
                      border: Border.all(color: AppColors.borderAlt),
                    ),
                    child: Row(
                      children: [
                        Container(
                          width: 6,
                          height: 6,
                          decoration: const BoxDecoration(
                            color: AppColors.primary,
                            shape: BoxShape.circle,
                          ),
                        ),
                        const SizedBox(width: 4),
                        const Text(
                          'live',
                          style: TextStyle(
                            color: AppColors.primary,
                            fontSize: 11,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),

            // ── Repo filter ──────────────────────────────────
            Container(
              decoration: const BoxDecoration(
                border: Border(bottom: BorderSide(color: AppColors.border)),
              ),
              child: Column(
                children: [
                  // Dropdown trigger
                  GestureDetector(
                    onTap: () =>
                        setState(() => _showRepoPicker = !_showRepoPicker),
                    child: Padding(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 16,
                        vertical: 10,
                      ),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Row(
                            children: [
                              const Icon(
                                Icons.account_tree_outlined,
                                size: 14,
                                color: AppColors.primary,
                              ),
                              const SizedBox(width: 8),
                              Text(
                                selectedRepoName,
                                style: const TextStyle(
                                  color: AppColors.primary,
                                  fontSize: 13,
                                  fontFamily: 'monospace',
                                ),
                              ),
                            ],
                          ),
                          Icon(
                            _showRepoPicker
                                ? Icons.keyboard_arrow_up
                                : Icons.keyboard_arrow_down,
                            size: 14,
                            color: AppColors.textSecondary,
                          ),
                        ],
                      ),
                    ),
                  ),

                  // Dropdown options
                  if (_showRepoPicker)
                    Container(
                      color: AppColors.surface,
                      child: Column(
                        children: [
                          // All repos option
                          _RepoOption(
                            label: 'All repos',
                            isSelected: selectedRepo == null,
                            onTap: () {
                              ref.read(selectedRepoIdProvider.notifier).state = null;
                              ref.read(currentPageProvider.notifier).state   = 1;
                              setState(() => _showRepoPicker = false);
                            },
                          ),
                          // Individual repos
                          ...repos.map((repo) => _RepoOption(
                            label: repo.repoFullName,
                            isSelected: selectedRepo == repo.repoId,
                            onTap: () {
                              ref.read(selectedRepoIdProvider.notifier).state =
                                  repo.repoId;
                              ref.read(currentPageProvider.notifier).state = 1;
                              setState(() => _showRepoPicker = false);
                            },
                          )),
                        ],
                      ),
                    ),
                ],
              ),
            ),

            // ── Events list ──────────────────────────────────
            Expanded(
              child: eventsAsync.when(
                // 'when' is like a switch on AsyncValue states
                // loading → show spinner
                // error   → show error message
                // data    → show the actual UI
                loading: () => const Center(
                  child: CircularProgressIndicator(color: AppColors.primary),
                ),
                error: (err, _) => Center(
                  child: Text(
                    'Error: $err',
                    style: const TextStyle(color: AppColors.error),
                  ),
                ),
                data: (response) {
                  final installed = hasInstallation.valueOrNull ?? true;
                  if (!installed) {
                    return Center(
                      child: Padding(
                        padding: const EdgeInsets.all(32),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            const Icon(Icons.integration_instructions_outlined, size: 48, color: AppColors.primary),
                            const SizedBox(height: 16),
                            const Text('Install GitHub App first',
                                style: TextStyle(color: AppColors.textPrimary, fontSize: 16, fontWeight: FontWeight.w600)),
                            const SizedBox(height: 8),
                            const Text('Connect your GitHub account to start monitoring repos',
                                style: TextStyle(color: AppColors.textSecondary, fontSize: 13),
                                textAlign: TextAlign.center),
                            const SizedBox(height: 24),
                            GestureDetector(
                              onTap: () async {
                                final uri = Uri.parse('https://github.com/apps/deploywatch/installations/new');
                                await launchUrl(uri, mode: LaunchMode.externalApplication);
                              },
                              child: Container(
                                padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                                decoration: BoxDecoration(
                                  color: AppColors.primaryDim,
                                  borderRadius: BorderRadius.circular(8),
                                  border: Border.all(color: AppColors.primary),
                                ),
                                child: const Text('Install GitHub App →',
                                    style: TextStyle(color: AppColors.primary, fontSize: 14)),
                              ),
                            ),
                          ],
                        ),
                      ),
                    );
                  }
                  if (response.events.isEmpty) {
                    return const Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(
                            Icons.flash_on_outlined,
                            size: 48,
                            color: AppColors.textDim,
                          ),
                          SizedBox(height: 16),
                          Text(
                            'No events yet',
                            style: TextStyle(color: AppColors.textSecondary),
                          ),
                        ],
                      ),
                    );
                  }

                  return RefreshIndicator(
                    color: AppColors.primary,
                    onRefresh: () => ref.refresh(eventsProvider.future),
                    child: ListView.builder(
                      itemCount: response.events.length + 1, // +1 for pagination
                      itemBuilder: (context, index) {
                        // Last item = pagination controls
                        if (index == response.events.length) {
                          return _Pagination(
                            currentPage: currentPage,
                            totalPages: response.totalPages,
                            onPrev: () => ref
                                .read(currentPageProvider.notifier)
                                .state--,
                            onNext: () => ref
                                .read(currentPageProvider.notifier)
                                .state++,
                          );
                        }

                        final event = response.events[index];
                        final repoName = repos
                            .where((r) => r.repoId == event.repoId)
                            .map((r) => r.repoFullName)
                            .firstOrNull ?? '';

                        return EventCard(
                          event: event,
                          repoFullName: repoName,
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
}

// ── Small helper widgets ───────────────────────────────────────

class _RepoOption extends StatelessWidget {
  final String label;
  final bool isSelected;
  final VoidCallback onTap;  // VoidCallback = () => void in TypeScript

  const _RepoOption({
    required this.label,
    required this.isSelected,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        decoration: const BoxDecoration(
          border: Border(bottom: BorderSide(color: AppColors.border)),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              label,
              style: TextStyle(
                color: isSelected ? AppColors.primary : AppColors.textSecondary,
                fontSize: 13,
                fontFamily: 'monospace',
              ),
            ),
            if (isSelected)
              const Icon(Icons.check, size: 14, color: AppColors.primary),
          ],
        ),
      ),
    );
  }
}

class _Pagination extends StatelessWidget {
  final int currentPage;
  final int totalPages;
  final VoidCallback onPrev;
  final VoidCallback onNext;

  const _Pagination({
    required this.currentPage,
    required this.totalPages,
    required this.onPrev,
    required this.onNext,
  });

  @override
  Widget build(BuildContext context) {
    if (totalPages <= 1) return const SizedBox.shrink(); // hide if only 1 page

    return Padding(
      padding: const EdgeInsets.all(16),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          GestureDetector(
            onTap: currentPage > 1 ? onPrev : null,
            child: Text(
              '← Prev',
              style: TextStyle(
                color: currentPage > 1
                    ? AppColors.primary
                    : AppColors.textDim,
              ),
            ),
          ),
          Text(
            '$currentPage / $totalPages',
            style: const TextStyle(
              color: AppColors.textSecondary,
              fontSize: 12,
            ),
          ),
          GestureDetector(
            onTap: currentPage < totalPages ? onNext : null,
            child: Text(
              'Next →',
              style: TextStyle(
                color: currentPage < totalPages
                    ? AppColors.primary
                    : AppColors.textDim,
              ),
            ),
          ),
        ],
      ),
    );
  }
}