import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:devlog_ai/screens/dashboard/dashboard_screen.dart';
import 'package:devlog_ai/screens/logs/logs_screen.dart';
import 'package:devlog_ai/screens/settings/settings_screen.dart';
import 'package:devlog_ai/services/auth_service.dart';

class AppDrawer extends ConsumerStatefulWidget {
  const AppDrawer({Key? key}) : super(key: key);

  @override
  ConsumerState<AppDrawer> createState() => _AppDrawerState();
}

class _AppDrawerState extends ConsumerState<AppDrawer> {
  String _userName = '';
  String _userEmail = '';

  @override
  void initState() {
    super.initState();
    _loadUserData();
  }

  // 載入使用者資料
  Future<void> _loadUserData() async {
    final userData = await ref.read(authServiceProvider).getUserData();
    if (mounted) {
      setState(() {
        _userName = userData['userName'] ?? '使用者';
        _userEmail = userData['userEmail'] ?? '';
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Drawer(
      child: Column(
        children: [
          // 抽屜標頭
          UserAccountsDrawerHeader(
            accountName: Text(_userName),
            accountEmail: Text(_userEmail),
            currentAccountPicture: CircleAvatar(
              backgroundColor: Theme.of(context).colorScheme.primary,
              child: Text(
                _userName.isNotEmpty ? _userName[0].toUpperCase() : 'U',
                style: const TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                  color: Colors.white,
                ),
              ),
            ),
            decoration: BoxDecoration(
              color: Theme.of(context).colorScheme.primary,
            ),
          ),
          
          // 選單項目
          ListTile(
            leading: const Icon(Icons.dashboard),
            title: const Text('儀表板'),
            onTap: () {
              Navigator.pushReplacement(
                context,
                MaterialPageRoute(builder: (context) => const DashboardScreen()),
              );
            },
          ),
          ListTile(
            leading: const Icon(Icons.list_alt),
            title: const Text('日誌'),
            onTap: () {
              Navigator.push(
                context,
                MaterialPageRoute(builder: (context) => const LogsScreen()),
              );
            },
          ),
          ListTile(
            leading: const Icon(Icons.analytics),
            title: const Text('分析'),
            onTap: () {
              // 導航到分析頁面
              Navigator.pop(context);
            },
          ),
          ListTile(
            leading: const Icon(Icons.settings),
            title: const Text('設定'),
            onTap: () {
              Navigator.push(
                context,
                MaterialPageRoute(builder: (context) => const SettingsScreen()),
              );
            },
          ),
          
          const Divider(),
          
          // 底部項目
          ListTile(
            leading: const Icon(Icons.help),
            title: const Text('說明'),
            onTap: () {
              // 顯示說明對話框
              Navigator.pop(context);
              _showHelpDialog(context);
            },
          ),
          ListTile(
            leading: const Icon(Icons.logout),
            title: const Text('登出'),
            onTap: () async {
              await ref.read(authStateProvider.notifier).logout();
              Navigator.pop(context);
            },
          ),
          
          const Spacer(),
          
          // 版本資訊
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: Text(
              'DevLog AI v1.0.0',
              style: TextStyle(
                color: Theme.of(context).colorScheme.onSurfaceVariant,
              ),
            ),
          ),
        ],
      ),
    );
  }

  // 顯示說明對話框
  void _showHelpDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: const Text('說明'),
          content: const SingleChildScrollView(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisSize: MainAxisSize.min,
              children: [
                Text('DevLog AI 是一個 AI 驅動的開發者日誌與事件分析平台。'),
                SizedBox(height: 16),
                Text('使用方法：'),
                SizedBox(height: 8),
                Text('1. 在儀表板中使用自然語言查詢您的日誌。'),
                Text('2. 查看系統狀態和最近的日誌。'),
                Text('3. 在日誌頁面中瀏覽和過濾所有日誌。'),
                Text('4. 在設定頁面中配置您的偏好設定。'),
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
          ],
        );
      },
    );
  }
}
