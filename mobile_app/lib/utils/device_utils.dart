import 'dart:convert';
import 'package:crypto/crypto.dart';
import 'package:uuid/uuid.dart';

class DeviceUtils {
  /// Generate a hashed device ID for privacy protection
  /// Uses UUID v4 to generate a unique identifier, then hashes it with SHA-256
  static Future<String> generateDeviceId() async {
    // Generate a UUID v4
    const uuid = Uuid();
    final rawId = uuid.v4();
    
    // Hash the UUID using SHA-256 for privacy protection
    final bytes = utf8.encode(rawId);
    final digest = sha256.convert(bytes);
    return digest.toString();
  }
}

