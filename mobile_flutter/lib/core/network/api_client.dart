import 'dart:ui';

import 'package:dio/dio.dart';
import '../config/app_config.dart';
import 'token_store.dart';

class ApiClient {
  late final Dio _dio;
  final TokenStore _tokenStore;
  VoidCallback? onUnauthorized;

  bool _isRefreshing = false;
  final List<_RetryRequest> _queue = [];

  ApiClient(this._tokenStore) {
    _dio = Dio(
      BaseOptions(
        baseUrl: AppConfig.baseURL,
        connectTimeout: const Duration(seconds: 10),
        receiveTimeout: const Duration(seconds: 30),
        headers: {'Content-Type': 'application/json'},
      ),
    );
    _dio.interceptors.add(_AuthInterceptor(this));
  }

  Future<T> get<T>(String path) async {
    final res = await _dio.get<Map<String, dynamic>>(path);
    return _unwrap<T>(res.data!);
  }

  Future<T> post<T>(String path, [Object? body]) async {
    final res = await _dio.post<Map<String, dynamic>>(path, data: body);
    return _unwrap<T>(res.data!);
  }

  Future<T> delete<T>(String path) async {
    final res = await _dio.delete<Map<String, dynamic>>(path);
    return _unwrap<T>(res.data!);
  }

  T _unwrap<T>(Map<String, dynamic> json) {
    final data = json['data'] ?? json;
    return data as T;
  }

  Future<String?> _getAccessToken() => _tokenStore.getAccessToken();

  Future<String?> refresh() async {
    final refreshToken = await _tokenStore.getRefreshToken();
    if (refreshToken == null) return null;
    try {
      final plainDio = Dio(BaseOptions(baseUrl: AppConfig.baseURL));
      final res = await plainDio.post<Map<String, dynamic>>(
        '/api/auth/refresh/mobile',
        data: {'refresh_token': refreshToken},
      );
      final data = res.data!['data'] as Map<String, dynamic>;
      final newAccess  = data['access_token']  as String;
      final newRefresh = data['refresh_token'] as String;
      await _tokenStore.setTokens(newAccess, newRefresh);
      return newAccess;
    } catch (_) {
      await _tokenStore.clear();
      return null;
    }
  }
}

class _AuthInterceptor extends Interceptor {
  final ApiClient _client;
  _AuthInterceptor(this._client);

  @override
  Future<void> onRequest(
      RequestOptions options,
      RequestInterceptorHandler handler,
      ) async {
    final token = await _client._getAccessToken();
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }

  @override
  Future<void> onError(
      DioException err,
      ErrorInterceptorHandler handler,
      ) async {
    if (err.response?.statusCode != 401) {
      handler.next(err);
      return;
    }
    final options = err.requestOptions;
    if (options.extra['_retried'] == true) {
      handler.next(err);
      return;
    }
    if (_client._isRefreshing) {
      _client._queue.add(_RetryRequest(options, handler));
      return;
    }
    _client._isRefreshing = true;
    final newToken = await _client.refresh();
    _client._isRefreshing = false;
    if (newToken == null) {
      for (final r in _client._queue) r.handler.next(err);
      _client._queue.clear();
      _client.onUnauthorized?.call(); // ← redirect to login
      handler.next(err);
      return;
    }
    for (final r in _client._queue) {
      r.options.headers['Authorization'] = 'Bearer $newToken';
      r.options.extra['_retried'] = true;
      try {
        final res = await _client._dio.fetch(r.options);
        r.handler.resolve(res);
      } catch (e) {
        r.handler.next(err);
      }
    }
    _client._queue.clear();
    options.headers['Authorization'] = 'Bearer $newToken';
    options.extra['_retried'] = true;
    try {
      final res = await _client._dio.fetch(options);
      handler.resolve(res);
    } catch (e) {
      handler.next(err);
    }
  }
}

class _RetryRequest {
  final RequestOptions options;
  final ErrorInterceptorHandler handler;
  _RetryRequest(this.options, this.handler);
}