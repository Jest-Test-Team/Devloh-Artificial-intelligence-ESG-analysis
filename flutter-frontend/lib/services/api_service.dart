import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:devlog_ai/models/log_entry.dart';
import 'package:devlog_ai/models/query_request.dart';
import 'package:devlog_ai/models/query_response.dart';
import 'package:devlog_ai/services/auth_service.dart';

// API 服務 Provider
final apiServiceProvider = Provider<ApiService>((ref) {
  final authService = ref.watch(authServiceProvider);
  return ApiService(authService: authService);
});

// API 服務
class ApiService {
  final AuthService authService;
  late final Dio _dio;
  
  // API 基礎 URL，可從環境變數取得
  final String baseUrl;
  
  ApiService({
    required this.authService,
    this.baseUrl = 'http://localhost:8080/api/v1',
  }) {
    _dio = Dio(BaseOptions(
      baseUrl: baseUrl,
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 10),
      contentType: 'application/json',
    ));
    
    // 添加攔截器，自動添加授權標頭
    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        final token = await authService.getAuthToken();
        if (token != null) {
          options.headers['Authorization'] = 'Bearer $token';
        }
        return handler.next(options);
      },
      onError: (DioException e, handler) {
        // 處理錯誤，例如 401 未授權
        if (e.response?.statusCode == 401) {
          // 可以在這裡觸發登出流程
        }
        return handler.next(e);
      },
    ));
  }
  
  // 健康檢查
  Future<bool> checkHealth() async {
    try {
      final response = await _dio.get('/health');
      return response.statusCode == 200;
    } catch (e) {
      return false;
    }
  }
  
  // 查詢日誌
  Future<QueryResponse> queryLogs(QueryRequest request) async {
    try {
      final response = await _dio.post(
        '/query',
        data: request.toJson(),
      );
      
      return QueryResponse.fromJson(response.data);
    } on DioException catch (e) {
      throw _handleDioError(e, 'Failed to query logs');
    } catch (e) {
      throw Exception('Unexpected error: $e');
    }
  }
  
  // 提交日誌
  Future<void> submitLog(LogEntry log) async {
    try {
      await _dio.post(
        '/logs',
        data: {
          'project_id': log.projectId,
          'logs': [log.toJson()],
        },
      );
    } on DioException catch (e) {
      throw _handleDioError(e, 'Failed to submit log');
    } catch (e) {
      throw Exception('Unexpected error: $e');
    }
  }
  
  // 批量提交日誌
  Future<void> submitLogs(String projectId, List<LogEntry> logs) async {
    try {
      await _dio.post(
        '/logs',
        data: {
          'project_id': projectId,
          'logs': logs.map((log) => log.toJson()).toList(),
        },
      );
    } on DioException catch (e) {
      throw _handleDioError(e, 'Failed to submit logs');
    } catch (e) {
      throw Exception('Unexpected error: $e');
    }
  }
  
  // 處理 Dio 錯誤
  Exception _handleDioError(DioException e, String defaultMessage) {
    if (e.response != null) {
      final statusCode = e.response!.statusCode;
      final responseData = e.response!.data;
      
      String errorMessage;
      if (responseData is Map && responseData.containsKey('message')) {
        errorMessage = responseData['message'];
      } else if (responseData is String && responseData.isNotEmpty) {
        errorMessage = responseData;
      } else {
        errorMessage = '$defaultMessage (Status: $statusCode)';
      }
      
      return Exception(errorMessage);
    }
    
    if (e.type == DioExceptionType.connectionTimeout) {
      return Exception('Connection timeout. Please try again.');
    } else if (e.type == DioExceptionType.receiveTimeout) {
      return Exception('Server took too long to respond. Please try again.');
    } else if (e.type == DioExceptionType.connectionError) {
      return Exception('No internet connection. Please check your network.');
    }
    
    return Exception('$defaultMessage: ${e.message}');
  }
}
