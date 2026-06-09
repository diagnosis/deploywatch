import 'dart:async';
import 'dart:convert';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:http/http.dart' as http;
import 'package:mobile_flutter/core/network/api_client.dart';
import '../config/app_config.dart';
import 'token_store.dart';
import 'providers.dart';

class SSEEvent {
  final String type;
  final String data;
  const SSEEvent({required this.type, required this.data});
}

class SSEService {
  final TokenStore _tokenStore;
  final ApiClient _apiClient;

  StreamController<SSEEvent>? _controller;
  http.Client? _client;
  bool _intentionallyClosed = false;
  Timer? _reconnectTimer;

  SSEService(this._tokenStore, this._apiClient);

  Stream<SSEEvent> get events {
    _controller ??= StreamController<SSEEvent>.broadcast();
    return _controller!.stream;
  }

  Future<void> connect() async {
    _intentionallyClosed = false;
    await _connect();
  }

  Future<void> _connect() async {
    _reconnectTimer?.cancel();

    final token = await _tokenStore.getAccessToken();
    if (token == null) {
      print('SSE: no token, aborting');
      return;
    }

    if (_controller == null || _controller!.isClosed) {
      _controller = StreamController<SSEEvent>.broadcast();
    }

    final url = Uri.parse('${AppConfig.baseURL}api/sse');
    print('SSE: connecting to $url');

    try {
      _client = http.Client();
      final request = http.Request('GET', url);
      request.headers['Authorization'] = 'Bearer $token';
      request.headers['Accept'] = 'text/event-stream';
      request.headers['Cache-Control'] = 'no-cache';

      final response = await _client!.send(request);

      if (response.statusCode != 200) {
        if (response.statusCode == 401) {
          print('SSE: 401, refreshing token...');
          final newToken = await _apiClient.refresh();
          if (newToken != null) {
            await _connect();
          } else {
            _intentionallyClosed = true; // logout handled by apiClient
          }
          return;
        }
        print('SSE: bad status ${response.statusCode}');
        _onDisconnect();
        return;
      }

      print('SSE: connected!');

      // SSE format:
      // event: push
      // data: {...json...}
      //
      // (blank line = end of event)
      String eventType = '';
      response.stream
          .transform(utf8.decoder)
          .transform(const LineSplitter())
          .listen(
            (line) {
          if (line.startsWith('event:')) {
            eventType = line.substring(6).trim();
          } else if (line.startsWith('data:')) {
            final data = line.substring(5).trim();
            if (eventType.isNotEmpty) {
              print('SSE: received event=$eventType');
              _controller?.add(SSEEvent(type: eventType, data: data));
              eventType = '';
            }
          }
          // blank line = event separator, ignore
        },
        onDone: _onDisconnect,
        onError: (_) => _onDisconnect(),
      );
    } catch (e) {
      print('SSE: connection failed: $e');
      _onDisconnect();
    }
  }

  void _onDisconnect() {
    if (_intentionallyClosed) return;
    print('SSE: disconnected, reconnecting in 3s...');
    _client?.close();
    _client = null;
    _reconnectTimer = Timer(const Duration(seconds: 3), _connect);
  }

  void disconnect() {
    _intentionallyClosed = true;
    _reconnectTimer?.cancel();
    _client?.close();
    _controller?.close();
    _client = null;
    _controller = null;
  }
}

final sseServiceProvider = Provider<SSEService>((ref) {
  final tokenStore = ref.watch(tokenStoreProvider);
  final apiClient = ref.watch(apiClientProvider);
  final service = SSEService(tokenStore, apiClient);
  ref.onDispose(service.disconnect);
  return service;
});