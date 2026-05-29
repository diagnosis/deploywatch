import 'dart:io';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/material.dart';
import '../core/network/api_client.dart';

class NotificationService {
  final ApiClient _apiClient;
  final FirebaseMessaging _messaging = FirebaseMessaging.instance;

  NotificationService(this._apiClient);

  Future<void> initialize() async {
    final settings = await _messaging.requestPermission(
      alert: true,
      badge: true,
      sound: true,
    );

    if (settings.authorizationStatus == AuthorizationStatus.denied) return;

    // Force delete cached token and get fresh one
    await _messaging.deleteToken();
    await Future.delayed(const Duration(seconds: 2));

    if (Platform.isIOS) {
      try {
        final fcmToken = await _messaging.getToken();
        print('FCM token: $fcmToken');
        if (fcmToken != null) await _registerToken(fcmToken);
      } catch (e) {
        print('FCM token error: $e');
      }
    }

    _messaging.onTokenRefresh.listen((token) async {
      print('FCM token refreshed: $token');
      await _registerToken(token);
    });

    FirebaseMessaging.onMessage.listen((message) {
      print('Foreground message: ${message.notification?.title}');
    });
  }
  Future<void> _registerToken(String token) async {
    try {
      await _apiClient.post<dynamic>('/api/device-tokens', {
        'token': token,
        'platform': 'ios',
      });
      print('Device token registered successfully');
    } catch (e) {
      print('Failed to register device token: $e');
    }
  }
}