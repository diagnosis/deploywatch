import 'package:flutter/foundation.dart';

// This is your Event type from types/index.ts
// In Dart we use classes instead of TypeScript interfaces
// 'fromJson' is the equivalent of JSON.parse — converts raw API response to a typed object

class Event {
  final String id;
  final int repoId;
  final String eventType;
  final String? action;      // nullable — not every event has an action
  final String actorLogin;
  final String payload;
  final String receivedAt;

  const Event({
    required this.id,
    required this.repoId,
    required this.eventType,
    this.action,
    required this.actorLogin,
    required this.payload,
    required this.receivedAt,
  });

  // Factory constructor — like a static method that returns an instance
  // This is how you deserialize JSON in Dart
  factory Event.fromJson(Map<String, dynamic> json) {
    return Event(
      id:          json['id']          as String,
      repoId:      json['repo_id']     as int,
      eventType:   json['event_type']  as String,
      action:      json['action']      as String?,
      actorLogin:  json['actor_login'] as String,
      payload:     json['payload']     as String,
      receivedAt:  json['received_at'] as String,
    );
  }
}

class EventsResponse {
  final List<Event> events;
  final int total;
  final int totalPages;
  final int page;
  final int limit;

  const EventsResponse({
    required this.events,
    required this.total,
    required this.totalPages,
    required this.page,
    required this.limit,
  });

  factory EventsResponse.fromJson(Map<String, dynamic> json) {
    // The API returns { events: [...], total: 10, ... }
    // We parse the events list by mapping each item through Event.fromJson
    final rawEvents = json['events'] as List<dynamic>? ?? [];
    return EventsResponse(
      events:     rawEvents.map((e) => Event.fromJson(e as Map<String, dynamic>)).toList(),
      total:      json['total']       as int? ?? 0,
      totalPages: json['total_pages'] as int? ?? 1,
      page:       json['page']        as int? ?? 1,
      limit:      json['limit']       as int? ?? 25,
    );
  }
}