import 'package:flutter_web_auth_2/flutter_web_auth_2.dart';
import '../../../core/network/token_store.dart';
import '../../../core/config/app_config.dart';

class AuthService {
  final TokenStore _tokenStore;

  AuthService(this._tokenStore);
  
  Future<bool> loginWithGithub() async{
    try{
      final result = await FlutterWebAuth2.authenticate(
        url: '${AppConfig.baseURL}/api/auth/github/login?mobile=true'
            '&redirect_uri=deploywatch://auth/callback',
        callbackUrlScheme: 'deploywatch', // must match your iOS URL scheme
      );
      final uri = Uri.parse(result);
      final accessToken = uri.queryParameters['access_token'];
      final refreshToken = uri.queryParameters['refresh_token'];

      if (accessToken == null || refreshToken == null) return false;

      await _tokenStore.setTokens(accessToken, refreshToken);
      return true;
    }catch(e){
      return false;
    }
  }

  Future<void> logout() async {
    await _tokenStore.clear();
  }
  
}