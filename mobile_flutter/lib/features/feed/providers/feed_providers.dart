import 'dart:async';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/network/providers.dart';
import '../../../core/network/sse_service.dart';  // ← changed
import '../data/event_service.dart';
import '../models/event_model.dart';

final eventsServiceProvider = Provider<EventService>((ref) {
  return EventService(ref.watch(apiClientProvider));
});

final selectedRepoIdProvider = StateProvider<int?>((ref) => null);
final currentPageProvider = StateProvider<int>((ref) => 1);

class EventsNotifier extends AsyncNotifier<EventsResponse> {
  @override
  Future<EventsResponse> build() async {
    final repoId = ref.watch(selectedRepoIdProvider);
    final page   = ref.watch(currentPageProvider);
    return ref.read(eventsServiceProvider).getEvents(
      repoId: repoId,
      page:   page,
    );
  }
}

final eventsProvider = AsyncNotifierProvider<EventsNotifier, EventsResponse>(
  EventsNotifier.new,
);

// SSE listener — replaces WebSocket
class SSENotifier extends AsyncNotifier<void> {
  StreamSubscription<SSEEvent>? _subscription;

  @override
  Future<void> build() async {
    final service = ref.watch(sseServiceProvider);
    await service.connect();

    _subscription = service.events.listen((_) {
      ref.invalidate(eventsProvider);
    });

    ref.onDispose(() => _subscription?.cancel());
  }
}

final sseNotifierProvider = AsyncNotifierProvider<SSENotifier, void>(
  SSENotifier.new,
);