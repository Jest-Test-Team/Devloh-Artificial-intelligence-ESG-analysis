import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';

// 模擬日誌資料
final mockLogsProvider = Provider<List<Map<String, dynamic>>>((ref) {
  return [
    {
      'timestamp': DateTime.now().subtract(const Duration(minutes: 5)),
      'level': 'ERROR',
      'message': '連接到資料庫失敗',
      'service': 'database-service',
    },
    {
      'timestamp': DateTime.now().subtract(const Duration(minutes: 10)),
      'level': 'WARN',
      'message': 'API 請求速率過高',
      'service': 'api-service',
    },
    {
      'timestamp': DateTime.now().subtract(const Duration(minutes: 15)),
      'level': 'INFO',
      'message': '使用者登入成功',
      'service': 'auth-service',
    },
    {
      'timestamp': DateTime.now().subtract(const Duration(minutes: 20)),
      'level': 'DEBUG',
      'message': '處理請求：GET /api/v1/users',
      'service': 'api-service',
    },
    {
      'timestamp': DateTime.now().subtract(const Duration(minutes: 25)),
      'level': 'INFO',
      'message': '服務啟動完成',
      'service': 'app-service',
    },
  ];
});

// 最近日誌卡片
class RecentLogsCard extends ConsumerWidget {
  final VoidCallback onViewAll;

  const RecentLogsCard({
    required this.onViewAll,
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final logs = ref.watch(mockLogsProvider);

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
                  '最近日誌',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                TextButton(
                  onPressed: onViewAll,
                  child: const Text('查看全部'),
                ),
              ],
            ),
            const SizedBox(height: 8),
            
            // 日誌列表
            ...logs.map((log) => _buildLogItem(context, log)).toList(),
          ],
        ),
      ),
    );
  }

  // 建立日誌項目
  Widget _buildLogItem(BuildContext context, Map<String, dynamic> log) {
    final timestamp = log['timestamp'] as DateTime;
    final level = log['level'] as String;
    final message = log['message'] as String;
    final service = log['service'] as String;
    
    // 根據日誌級別選擇顏色
    Color levelColor;
    switch (level) {
      case 'ERROR':
        levelColor = Colors.red;
        break;
      case 'WARN':
        levelColor = Colors.orange;
        break;
      case 'INFO':
        levelColor = Colors.blue;
        break;
      case 'DEBUG':
        levelColor = Colors.grey;
        break;
      default:
        levelColor = Colors.black;
    }
    
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 日誌級別標記
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
            decoration: BoxDecoration(
              color: levelColor.withOpacity(0.2),
              borderRadius: BorderRadius.circular(4),
            ),
            child: Text(
              level,
              style: TextStyle(
                color: levelColor,
                fontWeight: FontWeight.bold,
                fontSize: 12,
              ),
            ),
          ),
          const SizedBox(width: 8),
          
          // 日誌內容
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  message,
                  style: const TextStyle(
                    fontWeight: FontWeight.w500,
                  ),
                ),
                const SizedBox(height: 4),
                Row(
                  children: [
                    Text(
                      service,
                      style: TextStyle(
                        fontSize: 12,
                        color: Theme.of(context).colorScheme.secondary,
                      ),
                    ),
                    const SizedBox(width: 8),
                    Text(
                      DateFormat('HH:mm:ss').format(timestamp),
                      style: TextStyle(
                        fontSize: 12,
                        color: Theme.of(context).colorScheme.onSurfaceVariant,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
