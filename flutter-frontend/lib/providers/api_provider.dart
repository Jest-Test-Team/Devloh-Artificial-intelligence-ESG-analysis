import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../api/api_service.dart';

// API Service Provider
final apiServiceProvider = Provider<ApiService>((ref) {
  // Replace with actual API URL in production
  const apiUrl = 'http://localhost:8080';
  return ApiService(baseUrl: apiUrl);
});

// API Health Provider
final apiHealthProvider = FutureProvider<bool>((ref) async {
  final apiService = ref.watch(apiServiceProvider);
  return await apiService.checkHealth();
});
