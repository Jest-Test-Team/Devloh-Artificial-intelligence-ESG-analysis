import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:timeago/timeago.dart' as timeago;

class LogListItem extends StatelessWidget {
  final Map<String, dynamic> log;

  const LogListItem({
    required this.log,
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final timestamp = log['timestamp'] as DateTime;
    final level = log['level'] as String;
    final message = log['message'] as String;
    final service = log['service'] as String;
    final source = log['source'] as String;
    
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
    
    return InkWell(
      onTap: () {
        _showLogDetails(context);
      },
      borderRadius: BorderRadius.circular(8),
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 8.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 頂部行：時間戳和級別
            Row(
              children: [
                // 時間戳
                Text(
                  DateFormat('yyyy-MM-dd HH:mm:ss').format(timestamp),
                  style: TextStyle(
                    fontSize: 12,
                    color: Theme.of(context).colorScheme.onSurfaceVariant,
                  ),
                ),
                const SizedBox(width: 8),
                
                // 相對時間
                Text(
                  timeago.format(timestamp, locale: 'zh_TW'),
                  style: TextStyle(
                    fontSize: 12,
                    color: Theme.of(context).colorScheme.onSurfaceVariant,
                  ),
                ),
                const Spacer(),
                
                // 日誌級別
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
              ],
            ),
            const SizedBox(height: 8),
            
            // 日誌訊息
            Text(
              message,
              style: const TextStyle(
                fontWeight: FontWeight.w500,
              ),
            ),
            const SizedBox(height: 8),
            
            // 底部行：服務和來源
            Row(
              children: [
                // 服務
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: Theme.of(context).colorScheme.primary.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Text(
                    service,
                    style: TextStyle(
                      fontSize: 12,
                      color: Theme.of(context).colorScheme.primary,
                    ),
                  ),
                ),
                const SizedBox(width: 8),
                
                // 來源
                Expanded(
                  child: Text(
                    source,
                    style: TextStyle(
                      fontSize: 12,
                      color: Theme.of(context).colorScheme.onSurfaceVariant,
                    ),
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
  
  // 顯示日誌詳情對話框
  void _showLogDetails(BuildContext context) {
    final timestamp = log['timestamp'] as DateTime;
    final level = log['level'] as String;
    final message = log['message'] as String;
    final service = log['service'] as String;
    final source = log['source'] as String;
    
    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: const Text('日誌詳情'),
          content: SingleChildScrollView(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisSize: MainAxisSize.min,
              children: [
                _buildDetailRow(context, '時間', DateFormat('yyyy-MM-dd HH:mm:ss.SSS').format(timestamp)),
                _buildDetailRow(context, '級別', level),
                _buildDetailRow(context, '服務', service),
                _buildDetailRow(context, '來源', source),
                _buildDetailRow(context, '訊息', message),
                
                // 如果有更多資料，可以在這裡添加
                if (log.containsKey('metadata')) ...[
                  const Divider(),
                  const Text(
                    '元資料',
                    style: TextStyle(
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  _buildMetadata(context, log['metadata']),
                ],
              ],
            ),
          ),
          actions: [
            TextButton(
              onPressed: () {
                Navigator.pop(context);
              },
              child: const Text('關閉'),
            ),
            TextButton(
              onPressed: () {
                // 複製日誌到剪貼簿
                Navigator.pop(context);
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(
                    content: Text('已複製日誌到剪貼簿'),
                  ),
                );
              },
              child: const Text('複製'),
            ),
          ],
        );
      },
    );
  }
  
  // 建立詳情行
  Widget _buildDetailRow(BuildContext context, String label, String value) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            label,
            style: const TextStyle(
              fontWeight: FontWeight.bold,
              fontSize: 12,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            value,
            style: const TextStyle(
              fontFamily: 'monospace',
            ),
          ),
        ],
      ),
    );
  }
  
  // 建立元資料顯示
  Widget _buildMetadata(BuildContext context, Map<String, dynamic> metadata) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(8),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surfaceVariant,
        borderRadius: BorderRadius.circular(4),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: metadata.entries.map((entry) {
          return Padding(
            padding: const EdgeInsets.only(bottom: 4),
            child: Text(
              '${entry.key}: ${entry.value}',
              style: const TextStyle(
                fontFamily: 'monospace',
                fontSize: 12,
              ),
            ),
          );
        }).toList(),
      ),
    );
  }
}
