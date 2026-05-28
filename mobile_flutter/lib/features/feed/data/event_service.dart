

import 'package:mobile_flutter/core/network/api_client.dart';
import 'package:mobile_flutter/features/feed/models/event_model.dart';

class EventService {
  final ApiClient _client;
  EventService(this._client);

  Future<EventsResponse> getEvents({
    int? repoId,
    int page = 1,
    int limit = 25,
}) async {
    final params = {
      'page': page.toString(),
      'limit': limit.toString(),
      if (repoId != null) 'repo_id': repoId.toString()
    };
    final query = Uri(queryParameters: params).query;
    final data = await _client.get<Map<String, dynamic>>('api/events?$query');
    return EventsResponse.fromJson(data);
  }
}

