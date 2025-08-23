import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../api/api_service.dart';
import '../models/query_request.dart';
import '../models/query_response.dart';

// Query state
class QueryState {
  final bool isLoading;
  final String? error;
  final QueryResponse? response;
  final String currentQuery;

  QueryState({
    this.isLoading = false,
    this.error,
    this.response,
    this.currentQuery = '',
  });

  QueryState copyWith({
    bool? isLoading,
    String? error,
    QueryResponse? response,
    String? currentQuery,
  }) {
    return QueryState(
      isLoading: isLoading ?? this.isLoading,
      error: error,
      response: response ?? this.response,
      currentQuery: currentQuery ?? this.currentQuery,
    );
  }
}

// Query state notifier
class QueryNotifier extends StateNotifier<QueryState> {
  final ApiService _apiService;

  QueryNotifier(this._apiService) : super(QueryState());

  Future<void> submitQuery(String query, String? projectId) async {
    if (query.isEmpty) {
      state = state.copyWith(error: 'Query cannot be empty');
      return;
    }

    state = state.copyWith(isLoading: true, error: null, currentQuery: query);

    try {
      final request = QueryRequest(query: query, projectId: projectId);
      final response = await _apiService.queryLogs(request);
      state = state.copyWith(isLoading: false, response: response);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  void clearError() {
    if (state.error != null) {
      state = state.copyWith(error: null);
    }
  }

  void clearResponse() {
    state = QueryState(currentQuery: '');
  }
}

// Query provider
final queryProvider = StateNotifierProvider<QueryNotifier, QueryState>((ref) {
  final apiService = ref.watch(apiServiceProvider);
  return QueryNotifier(apiService);
});
