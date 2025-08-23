import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_markdown/flutter_markdown.dart';
import 'package:devlog_ai/models/query_request.dart';
import 'package:devlog_ai/models/query_response.dart';
import 'package:devlog_ai/screens/dashboard/dashboard_screen.dart';
import 'package:devlog_ai/services/api_service.dart';
import 'package:devlog_ai/services/auth_service.dart';

// 查詢結果 Provider
final queryResultProvider = StateNotifierProvider<QueryResultNotifier, AsyncValue<QueryResponse?>>((ref) {
  final apiService = ref.watch(apiServiceProvider);
  return QueryResultNotifier(apiService);
});

// 查詢結果狀態管理
class QueryResultNotifier extends StateNotifier<AsyncValue<QueryResponse?>> {
  final ApiService _apiService;

  QueryResultNotifier(this._apiService) : super(const AsyncValue.data(null));

  Future<void> executeQuery(String query, String? projectId) async {
    if (query.isEmpty) return;
    
    state = const AsyncValue.loading();
    
    try {
      final request = QueryRequest.withDefaultOptions(
        query: query,
        projectId: projectId,
      );
      
      final response = await _apiService.queryLogs(request);
      state = AsyncValue.data(response);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
    }
  }
}

// 查詢結果卡片
class LogQueryCard extends ConsumerWidget {
  const LogQueryCard({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // 監聽查詢輸入
    final queryInput = ref.watch(queryInputProvider);
    
    // 當查詢輸入變更時，執行查詢
    ref.listen<String>(queryInputProvider, (previous, next) async {
      if (next.isNotEmpty && next != previous) {
        final userData = await ref.read(authServiceProvider).getUserData();
        final projectId = userData['projectId'];
        await ref.read(queryResultProvider.notifier).executeQuery(next, projectId);
      }
    });
    
    // 監聽查詢結果
    final queryResult = ref.watch(queryResultProvider);
    
    return queryResult.when(
      data: (response) {
        if (response == null) {
          // 沒有查詢結果
          return const SizedBox.shrink();
        }
        
        // 顯示查詢結果
        return Card(
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text(
                      '查詢結果',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    Text(
                      '處理時間：${response.processingMs} ms',
                      style: const TextStyle(
                        fontSize: 12,
                        color: Colors.grey,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                
                // 回答
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: Theme.of(context).colorScheme.surfaceVariant,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: MarkdownBody(
                    data: response.answer,
                    styleSheet: MarkdownStyleSheet(
                      p: const TextStyle(fontSize: 16),
                    ),
                  ),
                ),
                
                // 相關日誌
                if (response.relatedLogs != null && response.relatedLogs!.isNotEmpty) ...[
                  const SizedBox(height: 16),
                  const Text(
                    '相關日誌',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Container(
                    width: double.infinity,
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: Theme.of(context).colorScheme.surface,
                      borderRadius: BorderRadius.circular(8),
                      border: Border.all(color: Theme.of(context).colorScheme.outline),
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: response.relatedLogs!.map((log) {
                        return Padding(
                          padding: const EdgeInsets.only(bottom: 8),
                          child: Text(
                            log,
                            style: const TextStyle(
                              fontFamily: 'monospace',
                              fontSize: 14,
                            ),
                          ),
                        );
                      }).toList(),
                    ),
                  ),
                ],
                
                // 來源
                if (response.sources != null && response.sources!.isNotEmpty) ...[
                  const SizedBox(height: 16),
                  const Text(
                    '來源',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Wrap(
                    spacing: 8,
                    runSpacing: 8,
                    children: response.sources!.map((source) {
                      return Chip(
                        label: Text(source),
                        backgroundColor: Theme.of(context).colorScheme.surfaceVariant,
                      );
                    }).toList(),
                  ),
                ],
              ],
            ),
          ),
        );
      },
      loading: () {
        // 載入中
        return Card(
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  '查詢結果',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 16),
                Center(
                  child: Column(
                    children: [
                      const CircularProgressIndicator(),
                      const SizedBox(height: 16),
                      Text(
                        '正在處理查詢：$queryInput',
                        textAlign: TextAlign.center,
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        );
      },
      error: (error, stackTrace) {
        // 錯誤
        return Card(
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  '查詢結果',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 16),
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: Theme.of(context).colorScheme.errorContainer,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text(
                        '處理查詢時發生錯誤',
                        style: TextStyle(
                          fontWeight: FontWeight.bold,
                          color: Colors.red,
                        ),
                      ),
                      const SizedBox(height: 8),
                      Text(error.toString()),
                      const SizedBox(height: 16),
                      ElevatedButton(
                        onPressed: () {
                          final userData = ref.read(authServiceProvider).getUserData();
                          userData.then((data) {
                            ref.read(queryResultProvider.notifier).executeQuery(
                              queryInput,
                              data['projectId'],
                            );
                          });
                        },
                        child: const Text('重試'),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }
}
