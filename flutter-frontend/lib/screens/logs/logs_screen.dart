import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:devlog_ai/screens/logs/widgets/log_filter_dialog.dart';
import 'package:devlog_ai/screens/logs/widgets/log_list_item.dart';
import 'package:devlog_ai/widgets/app_drawer.dart';
import 'package:intl/intl.dart';

// 模擬日誌資料
final logsProvider = Provider<List<Map<String, dynamic>>>((ref) {
  final now = DateTime.now();
  return [
    {
      'id': '1',
      'timestamp': now.subtract(const Duration(minutes: 5)),
      'level': 'ERROR',
      'message': '連接到資料庫失敗：Connection refused',
      'service': 'database-service',
      'source': 'DatabaseConnector.java:120',
    },
    {
      'id': '2',
      'timestamp': now.subtract(const Duration(minutes: 10)),
      'level': 'WARN',
      'message': 'API 請求速率過高：已達到限制 100/分鐘',
      'service': 'api-service',
      'source': 'RateLimiter.go:45',
    },
    {
      'id': '3',
      'timestamp': now.subtract(const Duration(minutes: 15)),
      'level': 'INFO',
      'message': '使用者 user_123 登入成功',
      'service': 'auth-service',
      'source': 'AuthController.js:78',
    },
    {
      'id': '4',
      'timestamp': now.subtract(const Duration(minutes: 20)),
      'level': 'DEBUG',
      'message': '處理請求：GET /api/v1/users?page=1&limit=10',
      'service': 'api-service',
      'source': 'ApiRouter.js:230',
    },
    {
      'id': '5',
      'timestamp': now.subtract(const Duration(minutes: 25)),
      'level': 'INFO',
      'message': '服務啟動完成，監聽端口 8080',
      'service': 'app-service',
      'source': 'Server.go:55',
    },
    {
      'id': '6',
      'timestamp': now.subtract(const Duration(minutes: 30)),
      'level': 'ERROR',
      'message': '無法處理付款：信用卡已過期',
      'service': 'payment-service',
      'source': 'PaymentProcessor.java:87',
    },
    {
      'id': '7',
      'timestamp': now.subtract(const Duration(minutes: 35)),
      'level': 'WARN',
      'message': '記憶體使用率高：85%',
      'service': 'monitoring-service',
      'source': 'ResourceMonitor.js:42',
    },
    {
      'id': '8',
      'timestamp': now.subtract(const Duration(minutes: 40)),
      'level': 'INFO',
      'message': '資料庫遷移成功：version 2.5.0',
      'service': 'database-service',
      'source': 'MigrationRunner.java:33',
    },
    {
      'id': '9',
      'timestamp': now.subtract(const Duration(minutes: 45)),
      'level': 'DEBUG',
      'message': '快取命中率：78.5%',
      'service': 'cache-service',
      'source': 'CacheStats.go:67',
    },
    {
      'id': '10',
      'timestamp': now.subtract(const Duration(minutes: 50)),
      'level': 'ERROR',
      'message': 'API 外部服務調用失敗：timeout after 30s',
      'service': 'api-service',
      'source': 'ExternalClient.js:112',
    },
  ];
});

// 過濾狀態 Provider
final logFilterProvider = StateNotifierProvider<LogFilterNotifier, LogFilter>((ref) {
  return LogFilterNotifier();
});

// 過濾狀態
class LogFilter {
  final String searchQuery;
  final Set<String> selectedLevels;
  final Set<String> selectedServices;
  final DateTime? startDate;
  final DateTime? endDate;

  LogFilter({
    this.searchQuery = '',
    Set<String>? selectedLevels,
    Set<String>? selectedServices,
    this.startDate,
    this.endDate,
  }) : 
    selectedLevels = selectedLevels ?? {'ERROR', 'WARN', 'INFO', 'DEBUG'},
    selectedServices = selectedServices ?? {};
  
  LogFilter copyWith({
    String? searchQuery,
    Set<String>? selectedLevels,
    Set<String>? selectedServices,
    DateTime? startDate,
    DateTime? endDate,
    bool? clearStartDate,
    bool? clearEndDate,
  }) {
    return LogFilter(
      searchQuery: searchQuery ?? this.searchQuery,
      selectedLevels: selectedLevels ?? this.selectedLevels,
      selectedServices: selectedServices ?? this.selectedServices,
      startDate: clearStartDate == true ? null : startDate ?? this.startDate,
      endDate: clearEndDate == true ? null : endDate ?? this.endDate,
    );
  }
}

// 過濾狀態管理
class LogFilterNotifier extends StateNotifier<LogFilter> {
  LogFilterNotifier() : super(LogFilter());
  
  void updateSearchQuery(String query) {
    state = state.copyWith(searchQuery: query);
  }
  
  void updateSelectedLevels(Set<String> levels) {
    state = state.copyWith(selectedLevels: levels);
  }
  
  void updateSelectedServices(Set<String> services) {
    state = state.copyWith(selectedServices: services);
  }
  
  void updateDateRange(DateTime? start, DateTime? end) {
    state = state.copyWith(
      startDate: start,
      endDate: end,
      clearStartDate: start == null,
      clearEndDate: end == null,
    );
  }
  
  void resetFilters() {
    state = LogFilter();
  }
}

// 過濾後的日誌 Provider
final filteredLogsProvider = Provider<List<Map<String, dynamic>>>((ref) {
  final logs = ref.watch(logsProvider);
  final filter = ref.watch(logFilterProvider);
  
  return logs.where((log) {
    // 過濾日誌級別
    if (!filter.selectedLevels.contains(log['level'])) {
      return false;
    }
    
    // 過濾服務
    if (filter.selectedServices.isNotEmpty && 
        !filter.selectedServices.contains(log['service'])) {
      return false;
    }
    
    // 過濾日期範圍
    final timestamp = log['timestamp'] as DateTime;
    if (filter.startDate != null && timestamp.isBefore(filter.startDate!)) {
      return false;
    }
    if (filter.endDate != null && timestamp.isAfter(filter.endDate!)) {
      return false;
    }
    
    // 過濾搜尋查詢
    if (filter.searchQuery.isNotEmpty) {
      final message = log['message'] as String;
      final service = log['service'] as String;
      final source = log['source'] as String;
      
      return message.toLowerCase().contains(filter.searchQuery.toLowerCase()) ||
             service.toLowerCase().contains(filter.searchQuery.toLowerCase()) ||
             source.toLowerCase().contains(filter.searchQuery.toLowerCase());
    }
    
    return true;
  }).toList();
});

// 可用服務列表 Provider
final availableServicesProvider = Provider<List<String>>((ref) {
  final logs = ref.watch(logsProvider);
  final services = logs.map((log) => log['service'] as String).toSet().toList();
  services.sort();
  return services;
});

// 日誌畫面
class LogsScreen extends ConsumerStatefulWidget {
  const LogsScreen({Key? key}) : super(key: key);

  @override
  ConsumerState<LogsScreen> createState() => _LogsScreenState();
}

class _LogsScreenState extends ConsumerState<LogsScreen> {
  final TextEditingController _searchController = TextEditingController();
  
  @override
  void initState() {
    super.initState();
    _searchController.addListener(() {
      ref.read(logFilterProvider.notifier).updateSearchQuery(_searchController.text);
    });
  }
  
  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }
  
  @override
  Widget build(BuildContext context) {
    final logs = ref.watch(filteredLogsProvider);
    final filter = ref.watch(logFilterProvider);
    final availableServices = ref.watch(availableServicesProvider);
    
    return Scaffold(
      appBar: AppBar(
        title: const Text('日誌'),
        actions: [
          // 過濾按鈕
          IconButton(
            icon: const Icon(Icons.filter_list),
            onPressed: () {
              _showFilterDialog(context, availableServices);
            },
          ),
          // 重置過濾器
          IconButton(
            icon: const Icon(Icons.clear_all),
            onPressed: () {
              ref.read(logFilterProvider.notifier).resetFilters();
              _searchController.clear();
            },
          ),
        ],
      ),
      drawer: const AppDrawer(),
      body: Column(
        children: [
          // 搜尋欄
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: TextField(
              controller: _searchController,
              decoration: InputDecoration(
                hintText: '搜尋日誌...',
                prefixIcon: const Icon(Icons.search),
                suffixIcon: _searchController.text.isNotEmpty
                    ? IconButton(
                        icon: const Icon(Icons.clear),
                        onPressed: () {
                          _searchController.clear();
                        },
                      )
                    : null,
                border: const OutlineInputBorder(),
              ),
            ),
          ),
          
          // 過濾器標籤
          if (_hasActiveFilters(filter)) ...[
            SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 16.0),
              child: Row(
                children: [
                  // 日誌級別過濾器
                  if (filter.selectedLevels.length < 4) ...[
                    Chip(
                      label: Text('級別: ${filter.selectedLevels.join(', ')}'),
                      onDeleted: () {
                        ref.read(logFilterProvider.notifier).updateSelectedLevels(
                          {'ERROR', 'WARN', 'INFO', 'DEBUG'},
                        );
                      },
                    ),
                    const SizedBox(width: 8),
                  ],
                  
                  // 服務過濾器
                  if (filter.selectedServices.isNotEmpty) ...[
                    Chip(
                      label: Text('服務: ${filter.selectedServices.join(', ')}'),
                      onDeleted: () {
                        ref.read(logFilterProvider.notifier).updateSelectedServices({});
                      },
                    ),
                    const SizedBox(width: 8),
                  ],
                  
                  // 日期範圍過濾器
                  if (filter.startDate != null || filter.endDate != null) ...[
                    Chip(
                      label: Text(
                        '日期: ${filter.startDate != null ? DateFormat('MM/dd HH:mm').format(filter.startDate!) : '開始'} - '
                        '${filter.endDate != null ? DateFormat('MM/dd HH:mm').format(filter.endDate!) : '結束'}',
                      ),
                      onDeleted: () {
                        ref.read(logFilterProvider.notifier).updateDateRange(null, null);
                      },
                    ),
                  ],
                ],
              ),
            ),
            const SizedBox(height: 8),
          ],
          
          // 日誌列表
          Expanded(
            child: logs.isEmpty
                ? const Center(
                    child: Text('沒有符合條件的日誌'),
                  )
                : ListView.separated(
                    padding: const EdgeInsets.all(16.0),
                    itemCount: logs.length,
                    separatorBuilder: (context, index) => const Divider(),
                    itemBuilder: (context, index) {
                      final log = logs[index];
                      return LogListItem(log: log);
                    },
                  ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          // 導航到日誌提交頁面
        },
        child: const Icon(Icons.add),
        tooltip: '新增日誌',
      ),
    );
  }
  
  // 顯示過濾對話框
  void _showFilterDialog(BuildContext context, List<String> availableServices) {
    showDialog(
      context: context,
      builder: (context) => LogFilterDialog(
        availableServices: availableServices,
      ),
    );
  }
  
  // 檢查是否有活動過濾器
  bool _hasActiveFilters(LogFilter filter) {
    return filter.selectedLevels.length < 4 ||
           filter.selectedServices.isNotEmpty ||
           filter.startDate != null ||
           filter.endDate != null;
  }
}
