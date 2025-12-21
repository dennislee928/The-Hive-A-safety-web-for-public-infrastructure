import 'package:shared_preferences/shared_preferences.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class StorageService {
  static const _storage = FlutterSecureStorage();
  static SharedPreferences? _prefs;
  
  static const String _deviceIdKey = 'device_id';
  static const String _apiKeyKey = 'api_key';
  static const String _currentZoneKey = 'current_zone';
  static const String _preferredLanguageKey = 'preferred_language';
  
  static Future<void> init() async {
    _prefs = await SharedPreferences.getInstance();
  }
  
  // Device ID (hashed, stored securely)
  static Future<String?> getDeviceId() async {
    return await _storage.read(key: _deviceIdKey);
  }
  
  static Future<void> saveDeviceId(String deviceId) async {
    await _storage.write(key: _deviceIdKey, value: deviceId);
  }
  
  // API Key (stored securely)
  static Future<String?> getApiKey() async {
    return await _storage.read(key: _apiKeyKey);
  }
  
  static Future<void> saveApiKey(String apiKey) async {
    await _storage.write(key: _apiKeyKey, value: apiKey);
  }
  
  // Current Zone (public storage is fine)
  static String? getCurrentZone() {
    return _prefs?.getString(_currentZoneKey);
  }
  
  static Future<void> saveCurrentZone(String zoneId) async {
    await _prefs?.setString(_currentZoneKey, zoneId);
  }
  
  // Preferred Language
  static String? getPreferredLanguage() {
    return _prefs?.getString(_preferredLanguageKey);
  }
  
  static Future<void> savePreferredLanguage(String language) async {
    await _prefs?.setString(_preferredLanguageKey, language);
  }
  
  // Clear all data (logout)
  static Future<void> clearAll() async {
    await _storage.deleteAll();
    await _prefs?.clear();
  }
}

