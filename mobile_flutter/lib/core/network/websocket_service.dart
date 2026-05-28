import 'dart:async';
import 'dart:io';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile_flutter/core/network/providers.dart';
import '../config/app_config.dart';
import 'token_store.dart';

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
  Timer? _pingTimer; // ← class field, not local variable

  WebSocketService(this._tokenStore);

  Stream<WSEvent> get events {
    _controller ??= StreamController<WSEvent>.broadcast();
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
      print('WS: no token, aborting');
      return;
    }

    // Recreate controller if closed
    if (_controller == null || _controller!.isClosed) {
      _controller = StreamController<WSEvent>.broadcast();
    }

    final wsUrl = AppConfig.baseURL
        .trimRight()
        .replaceAll(RegExp(r'/$'), '')
        .replaceFirst('https://', 'wss://')
        .replaceFirst('http://', 'ws://');

    print('WS: connecting to $wsUrl/api/ws');

    try {
      _socket = await WebSocket.connect(
        '$wsUrl/api/ws',
        headers: {'Authorization': 'Bearer $token'},
      );
      print('WS: connected!');

      // Start ping timer ONCE after connecting
      _pingTimer?.cancel();
      _pingTimer = Timer.periodic(const Duration(seconds: 30), (_) {
        if (_socket?.readyState == WebSocket.open) {
          _socket?.add('ping');
        }
      });

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
    _pingTimer?.cancel();
    _socket = null;
    _reconnectTimer = Timer(const Duration(seconds: 3), _connect);
  }

  void disconnect() {
    _intentionallyClosed = true;
    _pingTimer?.cancel();
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
  ref.onDispose(service.disconnect);
  return service;
});