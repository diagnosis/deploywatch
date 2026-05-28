import 'dart:async';
import 'dart:io';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile_flutter/core/network/providers.dart';
import '../config/app_config.dart';
import 'token_store.dart';

// WebSocket event — same shape as your Go Event struct
class WSEvent {
  final String type;
  final String data;

  const WSEvent({required this.type, required this.data});
}

class WebSocketService {
  final TokenStore _tokenStore;

  WebSocket? _socket;
  StreamController<WSEvent>? _controller;
  bool _intentionallyClosed = false;
  Timer? _reconnectTimer;

  WebSocketService(this._tokenStore);

  // The stream that the UI listens to
  // Stream = like an async generator in JS, pushes values over time
  Stream<WSEvent> get events {
    _controller ??= StreamController<WSEvent>.broadcast();
    return _controller!.stream;
  }

  Future<void> connect() async {
    _intentionallyClosed = false;
    await _connect();
  }

  Future<void> _connect() async {
    final token = await _tokenStore.getAccessToken();
    print(token);
    if (token == null) {
      print('WS: no token, aborting');
      return;
    }

    final wsUrl = AppConfig.baseURL
        .trimRight()
        .replaceAll(RegExp(r'/$'), '') // remove trailing slash
        .replaceFirst('https://', 'wss://')
        .replaceFirst('http://', 'ws://');

    print('WS: connecting to $wsUrl/api/ws');

    try {
      _socket = await WebSocket.connect(
        '$wsUrl/api/ws',
        headers: {'Authorization': 'Bearer $token'},
      );
      print('WS: connected!');

      _controller ??= StreamController<WSEvent>.broadcast();

      _socket!.listen(
            (data) {
          print('WS: received: $data');
          if (data is String) {
            final parts   = data.split('\n');
            final type    = parts.isNotEmpty ? parts[0] : 'unknown';
            final payload = parts.length > 1 ? parts.sublist(1).join('\n') : '';
            _controller?.add(WSEvent(type: type, data: payload));
          }
        },
        onDone: _onDisconnect,
        onError: (_) => _onDisconnect(),
      );
    } catch (e) {
      print('WS: connection failed: $e');
      _onDisconnect();
    }
  }

  void _onDisconnect() {
    if (_intentionallyClosed) return;
    // Auto reconnect after 3 seconds
    _reconnectTimer = Timer(const Duration(seconds: 3), _connect);
  }

  void disconnect() {
    _intentionallyClosed = true;
    _reconnectTimer?.cancel();
    _socket?.close();
    _controller?.close();
    _controller = null;
    _socket = null;
  }
}

final webSocketServiceProvider = Provider<WebSocketService>((ref) {
  final tokenStore = ref.watch(tokenStoreProvider);
  final service    = WebSocketService(tokenStore);

  // Automatically disconnect when the provider is disposed
  ref.onDispose(service.disconnect);

  return service;
});