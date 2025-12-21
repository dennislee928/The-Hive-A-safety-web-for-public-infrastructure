import 'dart:convert';
import 'dart:io';
import 'package:dio/dio.dart';
import 'package:erh_safety_app/models/api_response.dart';
import 'package:erh_safety_app/services/storage_service.dart';

class ApiService {
  static late Dio _dio;
  static const String baseUrl = 'http://localhost:8080/api/v1';
  static const String route2Base = '$baseUrl/route2';
  
  static void init() {
    _dio = Dio(BaseOptions(
      baseUrl: baseUrl,
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 10),
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
    ));
    
    // Add interceptors
    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        // Add device ID and API key to headers
        final deviceId = await StorageService.getDeviceId();
        final apiKey = await StorageService.getApiKey();
        
        if (deviceId != null && apiKey != null) {
          options.headers['X-Device-ID'] = deviceId;
          options.headers['X-API-Key'] = apiKey;
        }
        
        return handler.next(options);
      },
      onError: (error, handler) {
        // Handle API errors
        return handler.next(error);
      },
    ));
  }
  
  // Device Registration
  static Future<ApiResponse<Map<String, dynamic>>> registerDevice(String deviceId) async {
    try {
      final response = await _dio.post(
        '$route2Base/devices/register',
        data: {
          'device_id': deviceId,
          'platform': Platform.isIOS ? 'ios' : 'android',
          'app_version': '1.0.0', // TODO: Get from package info
        },
      );
      return ApiResponse.success(response.data);
    } on DioException catch (e) {
      return ApiResponse.error(_handleError(e));
    }
  }
  
  // Register Push Token
  static Future<ApiResponse<void>> registerPushToken(String deviceId, String pushToken) async {
    try {
      await _dio.post(
        '$route2Base/devices/$deviceId/push-token',
        data: {'push_token': pushToken},
      );
      return const ApiResponse.success(null);
    } on DioException catch (e) {
      return ApiResponse.error(_handleError(e));
    }
  }
  
  // Get Guidance
  static Future<ApiResponse<Map<String, dynamic>>> getGuidance({
    required String zoneId,
    required String currentZone,
    required String targetZone,
  }) async {
    try {
      final response = await _dio.get(
        '$route2Base/guidance',
        queryParameters: {
          'zone_id': zoneId,
          'current_zone': currentZone,
          'target_zone': targetZone,
        },
      );
      return ApiResponse.success(response.data);
    } on DioException catch (e) {
      return ApiResponse.error(_handleError(e));
    }
  }
  
  // Submit Report
  static Future<ApiResponse<Map<String, dynamic>>> submitReport({
    required String zoneId,
    required String content,
  }) async {
    try {
      final response = await _dio.post(
        '$baseUrl/reports',
        data: {
          'zone_id': zoneId,
          'content': content,
        },
      );
      return ApiResponse.success(response.data);
    } on DioException catch (e) {
      return ApiResponse.error(_handleError(e));
    }
  }
  
  // Request Assistance
  static Future<ApiResponse<Map<String, dynamic>>> requestAssistance({
    required String zoneId,
    String? subZone,
    required String requestType,
    required String urgency,
    String? description,
  }) async {
    try {
      final response = await _dio.post(
        '$route2Base/assistance',
        data: {
          'zone_id': zoneId,
          'sub_zone': subZone,
          'request_type': requestType,
          'urgency': urgency,
          'description': description,
        },
      );
      return ApiResponse.success(response.data);
    } on DioException catch (e) {
      return ApiResponse.error(_handleError(e));
    }
  }
  
  // Submit Feedback
  static Future<ApiResponse<Map<String, dynamic>>> submitFeedback({
    String? incidentId,
    required String guidanceClarity,
    required String guidanceTimeliness,
    String? suggestions,
  }) async {
    try {
      final response = await _dio.post(
        '$route2Base/feedback',
        data: {
          'incident_id': incidentId,
          'guidance_clarity': guidanceClarity,
          'guidance_timeliness': guidanceTimeliness,
          'suggestions': suggestions,
        },
      );
      return ApiResponse.success(response.data);
    } on DioException catch (e) {
      return ApiResponse.error(_handleError(e));
    }
  }
  
  // Get CAP Messages
  static Future<ApiResponse<List<dynamic>>> getCapMessages(String zoneId) async {
    try {
      final response = await _dio.get(
        '$baseUrl/cap/zone/$zoneId',
      );
      return ApiResponse.success(response.data);
    } on DioException catch (e) {
      return ApiResponse.error(_handleError(e));
    }
  }
  
  static String _handleError(DioException error) {
    if (error.response != null) {
      return error.response?.data['message'] ?? 'An error occurred';
    } else if (error.type == DioExceptionType.connectionTimeout) {
      return 'Connection timeout. Please check your internet connection.';
    } else if (error.type == DioExceptionType.receiveTimeout) {
      return 'Receive timeout. Please try again.';
    } else {
      return 'Network error. Please check your connection.';
    }
  }
}

