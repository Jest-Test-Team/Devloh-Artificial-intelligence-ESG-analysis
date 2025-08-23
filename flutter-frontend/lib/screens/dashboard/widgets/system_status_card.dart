import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fl_chart/fl_chart.dart';

// 模擬系統狀態資料
final systemStatusProvider = Provider<Map<String, dynamic>>((ref) {
  return {
    'services': [
      {'name': 'API Service', 'status': 'healthy'},
      {'name': 'Database', 'status': 'healthy'},
      {'name': 'Processor', 'status': 'healthy'},
      {'name': 'Ingestor', 'status': 'degraded'},
      {'name': 'Vector DB', 'status': 'healthy'},
    ],
    'metrics': {
      'logs_processed_today': 12583,
      'queries_today': 147,
      'avg_response_time': 245, // ms
    },
    'chart_data': [
      {'time': '00:00', 'value': 120},
      {'time': '04:00', 'value': 80},
      {'time': '08:00', 'value': 220},
      {'time': '12:00', 'value': 350},
      {'time': '16:00', 'value': 290},
      {'time': '20:00', 'value': 180},
    ],
  };
});

// 系統狀態卡片
class SystemStatusCard extends ConsumerWidget {
  const SystemStatusCard({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final status = ref.watch(systemStatusProvider);
    final services = status['services'] as List<dynamic>;
    final metrics = status['metrics'] as Map<String, dynamic>;
    final chartData = status['chart_data'] as List<dynamic>;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              '系統狀態',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16),
            
            // 服務狀態
            const Text(
              '服務狀態',
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            ...services.map((service) => _buildServiceStatus(context, service)).toList(),
            const SizedBox(height: 16),
            
            // 指標
            const Text(
              '今日指標',
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            _buildMetricItem(
              context,
              '處理的日誌',
              metrics['logs_processed_today'].toString(),
              Icons.analytics,
            ),
            _buildMetricItem(
              context,
              '查詢次數',
              metrics['queries_today'].toString(),
              Icons.search,
            ),
            _buildMetricItem(
              context,
              '平均響應時間',
              '${metrics['avg_response_time']} ms',
              Icons.timer,
            ),
            const SizedBox(height: 16),
            
            // 圖表
            const Text(
              '日誌活動',
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            SizedBox(
              height: 150,
              child: _buildActivityChart(context, chartData),
            ),
          ],
        ),
      ),
    );
  }

  // 建立服務狀態項目
  Widget _buildServiceStatus(BuildContext context, Map<String, dynamic> service) {
    final name = service['name'] as String;
    final status = service['status'] as String;
    
    Color statusColor;
    switch (status) {
      case 'healthy':
        statusColor = Colors.green;
        break;
      case 'degraded':
        statusColor = Colors.orange;
        break;
      case 'down':
        statusColor = Colors.red;
        break;
      default:
        statusColor = Colors.grey;
    }
    
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4.0),
      child: Row(
        children: [
          Container(
            width: 12,
            height: 12,
            decoration: BoxDecoration(
              color: statusColor,
              shape: BoxShape.circle,
            ),
          ),
          const SizedBox(width: 8),
          Text(name),
          const Spacer(),
          Text(
            status.toUpperCase(),
            style: TextStyle(
              fontSize: 12,
              fontWeight: FontWeight.bold,
              color: statusColor,
            ),
          ),
        ],
      ),
    );
  }

  // 建立指標項目
  Widget _buildMetricItem(
    BuildContext context,
    String label,
    String value,
    IconData icon,
  ) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4.0),
      child: Row(
        children: [
          Icon(
            icon,
            size: 16,
            color: Theme.of(context).colorScheme.primary,
          ),
          const SizedBox(width: 8),
          Text(label),
          const Spacer(),
          Text(
            value,
            style: const TextStyle(
              fontWeight: FontWeight.bold,
            ),
          ),
        ],
      ),
    );
  }

  // 建立活動圖表
  Widget _buildActivityChart(BuildContext context, List<dynamic> chartData) {
    final spots = chartData.asMap().entries.map((entry) {
      final index = entry.key.toDouble();
      final item = entry.value as Map<String, dynamic>;
      final value = (item['value'] as int).toDouble();
      return FlSpot(index, value);
    }).toList();
    
    return LineChart(
      LineChartData(
        gridData: FlGridData(show: false),
        titlesData: FlTitlesData(
          bottomTitles: AxisTitles(
            sideTitles: SideTitles(
              showTitles: true,
              getTitlesWidget: (value, meta) {
                if (value.toInt() >= 0 && value.toInt() < chartData.length) {
                  final item = chartData[value.toInt()] as Map<String, dynamic>;
                  return Padding(
                    padding: const EdgeInsets.only(top: 8.0),
                    child: Text(
                      item['time'] as String,
                      style: const TextStyle(
                        fontSize: 10,
                        color: Colors.grey,
                      ),
                    ),
                  );
                }
                return const Text('');
              },
              reservedSize: 30,
            ),
          ),
          leftTitles: AxisTitles(
            sideTitles: SideTitles(showTitles: false),
          ),
          rightTitles: AxisTitles(
            sideTitles: SideTitles(showTitles: false),
          ),
          topTitles: AxisTitles(
            sideTitles: SideTitles(showTitles: false),
          ),
        ),
        borderData: FlBorderData(show: false),
        lineBarsData: [
          LineChartBarData(
            spots: spots,
            isCurved: true,
            color: Theme.of(context).colorScheme.primary,
            barWidth: 3,
            isStrokeCapRound: true,
            dotData: FlDotData(show: false),
            belowBarData: BarAreaData(
              show: true,
              color: Theme.of(context).colorScheme.primary.withOpacity(0.2),
            ),
          ),
        ],
      ),
    );
  }
}
