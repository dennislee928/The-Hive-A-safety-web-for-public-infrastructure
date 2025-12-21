import 'package:flutter/foundation.dart';
import 'package:erh_safety_app/services/api_service.dart';

class GuidanceProvider with ChangeNotifier {
  Map<String, dynamic>? _currentGuidance;
  bool _isLoading = false;
  String? _error;
  
  Map<String, dynamic>? get currentGuidance => _currentGuidance;
  bool get isLoading => _isLoading;
  String? get error => _error;
  
  Future<bool> fetchGuidance({
    required String zoneId,
    required String currentZone,
    required String targetZone,
  }) async {
    _isLoading = true;
    _error = null;
    notifyListeners();
    
    try {
      final response = await ApiService.getGuidance(
        zoneId: zoneId,
        currentZone: currentZone,
        targetZone: targetZone,
      );
      
      if (response.success && response.data != null) {
        _currentGuidance = response.data!['guidance'] as Map<String, dynamic>;
        _isLoading = false;
        notifyListeners();
        return true;
      } else {
        _error = response.error ?? 'Failed to fetch guidance';
        _isLoading = false;
        notifyListeners();
        return false;
      }
    } catch (e) {
      _error = 'Error fetching guidance: $e';
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }
  
  void clearGuidance() {
    _currentGuidance = null;
    _error = null;
    notifyListeners();
  }
  
  void clearError() {
    _error = null;
    notifyListeners();
  }
}

