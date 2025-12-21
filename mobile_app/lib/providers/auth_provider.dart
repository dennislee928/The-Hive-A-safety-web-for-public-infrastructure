import 'package:flutter/foundation.dart';
import 'package:erh_safety_app/services/api_service.dart';
import 'package:erh_safety_app/services/storage_service.dart';
import 'package:erh_safety_app/utils/device_utils.dart';

class AuthProvider with ChangeNotifier {
  bool _isAuthenticated = false;
  bool _isLoading = false;
  String? _deviceId;
  String? _apiKey;
  String? _error;
  
  bool get isAuthenticated => _isAuthenticated;
  bool get isLoading => _isLoading;
  String? get deviceId => _deviceId;
  String? get error => _error;
  
  AuthProvider() {
    _checkAuthStatus();
  }
  
  Future<void> _checkAuthStatus() async {
    final savedDeviceId = await StorageService.getDeviceId();
    final savedApiKey = await StorageService.getApiKey();
    
    if (savedDeviceId != null && savedApiKey != null) {
      _deviceId = savedDeviceId;
      _apiKey = savedApiKey;
      _isAuthenticated = true;
      notifyListeners();
    }
  }
  
  Future<bool> registerDevice() async {
    _isLoading = true;
    _error = null;
    notifyListeners();
    
    try {
      // Generate device ID (hashed)
      final deviceId = await DeviceUtils.generateDeviceId();
      
      // Register with backend
      final response = await ApiService.registerDevice(deviceId);
      
      if (response.success && response.data != null) {
        _deviceId = deviceId;
        // API key can be in 'api_key' or 'device.api_key'
        String? apiKey;
        if (response.data!.containsKey('api_key')) {
          apiKey = response.data!['api_key'] as String?;
        } else if (response.data!.containsKey('device')) {
          final device = response.data!['device'] as Map<String, dynamic>?;
          apiKey = device?['api_key'] as String?;
        }
        
        if (apiKey != null) {
          _apiKey = apiKey;
          _isAuthenticated = true;
        
        // Save to storage
        await StorageService.saveDeviceId(deviceId);
        await StorageService.saveApiKey(_apiKey!);
        
        _isLoading = false;
        notifyListeners();
        return true;
      } else {
        _error = response.error ?? 'Registration failed';
        _isLoading = false;
        notifyListeners();
        return false;
      }
    } catch (e) {
      _error = 'Registration error: $e';
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }
  
  Future<void> logout() async {
    await StorageService.clearAll();
    _isAuthenticated = false;
    _deviceId = null;
    _apiKey = null;
    _error = null;
    notifyListeners();
  }
  
  void clearError() {
    _error = null;
    notifyListeners();
  }
}

