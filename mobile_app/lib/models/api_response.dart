class ApiResponse<T> {
  final bool success;
  final T? data;
  final String? error;
  
  const ApiResponse({
    required this.success,
    this.data,
    this.error,
  });
  
  factory ApiResponse.success(T data) {
    return ApiResponse(success: true, data: data);
  }
  
  factory ApiResponse.error(String error) {
    return ApiResponse(success: false, error: error);
  }
}

