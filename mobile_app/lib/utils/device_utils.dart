import 'dart:io';
import 'dart:convert';
import 'package:crypto/crypto.dart';
import 'package:device_info_plus/device_info_plus.dart';
import 'package:platform_device_id/platform_device_id.dart';

class DeviceUtils {
  /// Generate a hashed device ID for privacy protection
  static Future<String> generateDeviceId() async {
    try {
      String? deviceId;
      
      // Try to get platform-specific device ID
      if (Platform.isAndroid) {
        final androidInfo = await DeviceInfoPlugin().androidInfo;
        deviceId = androidInfo.id; // Android ID
      } else if (Platform.isIOS) {
        deviceId = await PlatformDeviceId.getDeviceId;
      }
      
      // If no device ID available, generate a UUID-based one
      if (deviceId == null || deviceId.isEmpty) {
        deviceId = DateTime.now().millisecondsSinceEpoch.toString();
      }
      
      // Hash the device ID using SHA-256
      final bytes = utf8.encode(deviceId);
      final digest = sha256.convert(bytes);
      return digest.toString();
    } catch (e) {
      // Fallback: use timestamp-based hash
      final timestamp = DateTime.now().millisecondsSinceEpoch.toString();
      final bytes = utf8.encode(timestamp);
      final digest = sha256.convert(bytes);
      return digest.toString();
    }
  }
}

