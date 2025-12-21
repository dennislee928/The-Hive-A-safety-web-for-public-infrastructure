import 'package:flutter/foundation.dart';
import 'package:erh_safety_app/services/api_service.dart';

class CapProvider with ChangeNotifier {
  List<Map<String, dynamic>> _capMessages = [];
  bool _isLoading = false;
  String? _error;
  
  List<Map<String, dynamic>> get capMessages => _capMessages;
  bool get isLoading => _isLoading;
  String? get error => _error;
  
  Future<bool> fetchCapMessages(String zoneId) async {
    _isLoading = true;
    _error = null;
    notifyListeners();
    
    try {
      final response = await ApiService.getCapMessages(zoneId);
      
      if (response.success && response.data != null) {
        _capMessages = List<Map<String, dynamic>>.from(response.data!);
        _isLoading = false;
        notifyListeners();
        return true;
      } else {
        _error = response.error ?? 'Failed to fetch CAP messages';
        _isLoading = false;
        notifyListeners();
        return false;
      }
    } catch (e) {
      _error = 'Error fetching CAP messages: $e';
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }
  
  void clearCapMessages() {
    _capMessages = [];
    _error = null;
    notifyListeners();
  }
  
  void clearError() {
    _error = null;
    notifyListeners();
  }
}

