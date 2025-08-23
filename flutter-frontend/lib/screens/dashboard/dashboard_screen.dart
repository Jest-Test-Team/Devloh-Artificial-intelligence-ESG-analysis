import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:devlog_ai/models/query_request.dart';
import 'package:devlog_ai/screens/dashboard/widgets/log_query_card.dart';
import 'package:devlog_ai/screens/dashboard/widgets/recent_logs_card.dart';
import 'package:devlog_ai/screens/dashboard/widgets/system_status_card.dart';
import 'package:devlog_ai/screens/logs/logs_screen.dart';
import 'package:devlog_ai/services/auth_service.dart';
import 'package:devlog_ai/widgets/app_drawer.dart';

// 查詢輸入 Provider
final queryInputProvider = StateProvider<String>((ref) => '');

// 儀表板畫面
class DashboardScreen extends ConsumerStatefulWidget {
  const DashboardScreen({Key? key}) : super(key: key);

  @override
  ConsumerState<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends ConsumerState<DashboardScreen> {
  final TextEditingController _queryController = TextEditingController();
  String _userName = '';
  
  @override
  void initState() {
    super.initState();
    _loadUserData();
  }
  
  @override
  void dispose() {
    _queryController.dispose();
    super.dispose();
  }
  
  // 載入使用者資料
  Future<void> _loadUserData() async {
    final userData = await ref.read(authServiceProvider).getUserData();
    if (mounted) {
      setState(() {
        _userName = userData['userName'] ?? '使用者';
      });
    }
  }
  
  // 登出
  Future<void> _logout() async {
    await ref.read(authStateProvider.notifier).logout();
  }
  
  // 提交查詢
  void _submitQuery() {
    final query = _queryController.text.trim();
    if (query.isEmpty) return;
    
    // 更新查詢輸入 Provider
    ref.read(queryInputProvider.notifier).state = query;
    
    // 清空輸入框
    _queryController.clear();
    
    // 隱藏鍵盤
    FocusScope.of(context).unfocus();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('DevLog AI 儀表板'),
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: _logout,
            tooltip: '登出',
          ),
        ],
      ),
      drawer: const AppDrawer(),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // 歡迎訊息
              Text(
                '歡迎，$_userName',
                style: Theme.of(context).textTheme.headlineSmall,
              ),
              const SizedBox(height: 24),
              
              // 查詢卡片
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text(
                        '用自然語言查詢您的日誌',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 16),
                      TextField(
                        controller: _queryController,
                        decoration: InputDecoration(
                          hintText: '例如：顯示過去一小時內的所有錯誤日誌',
                          border: const OutlineInputBorder(),
                          suffixIcon: IconButton(
                            icon: const Icon(Icons.search),
                            onPressed: _submitQuery,
                          ),
                        ),
                        onSubmitted: (_) => _submitQuery(),
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 24),
              
              // 查詢結果卡片
              const LogQueryCard(),
              const SizedBox(height: 24),
              
              // 系統狀態和最近日誌
              Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // 系統狀態卡片
                  const Expanded(
                    flex: 1,
                    child: SystemStatusCard(),
                  ),
                  const SizedBox(width: 16),
                  
                  // 最近日誌卡片
                  Expanded(
                    flex: 2,
                    child: RecentLogsCard(
                      onViewAll: () {
                        Navigator.push(
                          context,
                          MaterialPageRoute(
                            builder: (context) => const LogsScreen(),
                          ),
                        );
                      },
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
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
}