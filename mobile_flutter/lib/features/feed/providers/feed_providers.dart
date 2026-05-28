import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile_flutter/core/network/providers.dart';
import 'package:mobile_flutter/features/feed/data/event_service.dart';
import 'package:mobile_flutter/features/feed/models/event_model.dart';

import '../../../core/network/websocket_service.dart';

// Service provide


final eventsServiceProvider = Provider<EventService>((ref) {
  return EventService(ref.watch(apiClientProvider));
});

final selectedRepoIdProvider = StateProvider<int?>((ref) => null);

final currentPageProvider = StateProvider<int>((ref) => 1);

class EventsNotifier extends AsyncNotifier<EventsResponse> {
  @override
  Future<EventsResponse> build() async{
    final repoId = ref.watch(selectedRepoIdProvider);
    final page = ref.watch(currentPageProvider);

    return ref.read(eventsServiceProvider).getEvents(
      repoId: repoId,
      page: page,
    );
  }
}

final eventsProvider = AsyncNotifierProvider<EventsNotifier, EventsResponse>(
  EventsNotifier.new,
);

// WebSocket listener — connects and refreshes feed on new events
class WebSocketNotifier extends AsyncNotifier<void> {
  StreamSubscription<WSEvent>? _subscription;

  @override
  Future<void> build() async {
    final service = ref.watch(webSocketServiceProvider);
    await service.connect();

    _subscription = service.events.listen((_) {
      // New event received — refresh the feed
      ref.invalidateSelf();
      ref.invalidate(eventsProvider);
    });

    ref.onDispose(() => _subscription?.cancel());
  }
}

final webSocketNotifierProvider =
AsyncNotifierProvider<WebSocketNotifier, void>(
  WebSocketNotifier.new,
);