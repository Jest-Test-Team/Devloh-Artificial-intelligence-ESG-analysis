import 'package:dio/dio.dart';
import '../models/query_request.dart';
import '../models/query_response.dart';
import '../models/log_entry.dart';

class ApiService {
  final Dio _dio;
  final String baseUrl;

  ApiService({required this.baseUrl})
      : _dio = Dio(BaseOptions(
          baseUrl: baseUrl,
          connectTimeout: const Duration(seconds: 10),
          receiveTimeout: const Duration(seconds: 10),
          headers: {
            'Content-Type': 'application/json',
          },
        ));

  // Health check
  Future<bool> checkHealth() async {
    try {
      final response = await _dio.get('/health');
      return response.statusCode == 200;
    } catch (e) {
      return false;
    }
  }

  // Query logs
  Future<QueryResponse> queryLogs(QueryRequest request) async {
    try {
      final response = await _dio.post(
        '/api/v1/query',
        data: request.toJson(),
      );
      return QueryResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  // Send logs
  Future<void> sendLogs(String projectId, List<LogEntry> logs) async {
    try {
      await _dio.post(
        '/api/v1/logs',
        data: {
          'project_id': projectId,
          'logs': logs.map((log) => log.toJson()).toList(),
        },
      );
    } catch (e) {
      throw _handleError(e);
    }
  }

  // Handle Dio errors
  Exception _handleError(dynamic error) {
    if (error is DioException) {
      if (error.response != null) {
        return Exception(
            'API Error: ${error.response?.statusCode} - ${error.response?.data}');
      }
      return Exception('Network Error: ${error.message}');
    }
    return Exception('Unexpected Error: $error');
  }
}
