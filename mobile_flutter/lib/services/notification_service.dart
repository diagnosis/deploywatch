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

    print('Notification permission: ${settings.authorizationStatus}');
    if (settings.authorizationStatus == AuthorizationStatus.denied) return;

    if (Platform.isIOS) {
      // Retry up to 10 times with 2 second intervals
      for (int i = 0; i < 10; i++) {
        await Future.delayed(const Duration(seconds: 2));
        try {
          final fcmToken = await _messaging.getToken();
          print('FCM token: $fcmToken');
          if (fcmToken != null) {
            await _registerToken(fcmToken);
            break; // success, stop retrying
          }
        } catch (e) {
          print('Attempt ${i+1} failed: $e');
        }
      }
    }

    // This fires when FCM token is ready — catches cases where APNs token arrives late
    _messaging.onTokenRefresh.listen((token) async {
      print('FCM token via refresh: $token');
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