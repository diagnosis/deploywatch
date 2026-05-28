import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class TokenStore {
  final _storage = const FlutterSecureStorage(
    iOptions: IOSOptions(accessibility: KeychainAccessibility.first_unlock),
  );

  static const _accessKey  = 'access_token';
  static const _refreshKey = 'refresh_token';

  Future<String?> getAccessToken()  => _storage.read(key: _accessKey);
  Future<String?> getRefreshToken() => _storage.read(key: _refreshKey);

  Future<void> setTokens(String access, String refresh) async {
    await Future.wait([
      _storage.write(key: _accessKey,  value: access),
      _storage.write(key: _refreshKey, value: refresh),
    ]);
  }

  Future<void> clear() async {
    await Future.wait([
      _storage.delete(key: _accessKey),
      _storage.delete(key: _refreshKey),
    ]);
  }

  Future<bool> hasToken() async {
    final token = await getAccessToken();
    return token != null && token.isNotEmpty;
  }
}