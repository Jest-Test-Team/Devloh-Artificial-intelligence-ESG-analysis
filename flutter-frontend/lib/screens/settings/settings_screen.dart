import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:devlog_ai/widgets/app_drawer.dart';
import 'package:shared_preferences/shared_preferences.dart';

// 設定 Provider
final settingsProvider = StateNotifierProvider<SettingsNotifier, Map<String, dynamic>>((ref) {
  return SettingsNotifier();
});

// 設定狀態管理
class SettingsNotifier extends StateNotifier<Map<String, dynamic>> {
  SettingsNotifier() : super({
    'theme': 'system',
    'apiUrl': 'http://localhost:8080/api/v1',
    'logRetention': 30,
    'showDebugLogs': true,
    'enableNotifications': true,
  }) {
    _loadSettings();
  }
  
  // 載入設定
  Future<void> _loadSettings() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      
      final theme = prefs.getString('theme') ?? 'system';
      final apiUrl = prefs.getString('apiUrl') ?? 'http://localhost:8080/api/v1';
      final logRetention = prefs.getInt('logRetention') ?? 30;
      final showDebugLogs = prefs.getBool('showDebugLogs') ?? true;
      final enableNotifications = prefs.getBool('enableNotifications') ?? true;
      
      state = {
        'theme': theme,
        'apiUrl': apiUrl,
        'logRetention': logRetention,
        'showDebugLogs': showDebugLogs,
        'enableNotifications': enableNotifications,
      };
    } catch (e) {
      // 如果載入失敗，使用預設值
    }
  }
  
  // 更新設定
  Future<void> updateSetting(String key, dynamic value) async {
    try {
      final prefs = await SharedPreferences.getInstance();
      
      switch (key) {
        case 'theme':
          await prefs.setString(key, value);
          break;
        case 'apiUrl':
          await prefs.setString(key, value);
          break;
        case 'logRetention':
          await prefs.setInt(key, value);
          break;
        case 'showDebugLogs':
          await prefs.setBool(key, value);
          break;
        case 'enableNotifications':
          await prefs.setBool(key, value);
          break;
      }
      
      state = {...state, key: value};
    } catch (e) {
      // 如果儲存失敗，不更新狀態
    }
  }
}

// 設定畫面
class SettingsScreen extends ConsumerWidget {
  const SettingsScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final settings = ref.watch(settingsProvider);
    
    return Scaffold(
      appBar: AppBar(
        title: const Text('設定'),
      ),
      drawer: const AppDrawer(),
      body: ListView(
        padding: const EdgeInsets.all(16.0),
        children: [
          // 外觀設定
          _buildSectionHeader(context, '外觀'),
          _buildThemeSetting(context, ref, settings['theme']),
          const Divider(),
          
          // API 設定
          _buildSectionHeader(context, 'API 設定'),
          _buildApiUrlSetting(context, ref, settings['apiUrl']),
          const Divider(),
          
          // 日誌設定
          _buildSectionHeader(context, '日誌設定'),
          _buildLogRetentionSetting(context, ref, settings['logRetention']),
          _buildShowDebugLogsSetting(context, ref, settings['showDebugLogs']),
          const Divider(),
          
          // 通知設定
          _buildSectionHeader(context, '通知設定'),
          _buildEnableNotificationsSetting(context, ref, settings['enableNotifications']),
          const Divider(),
          
          // 關於
          _buildSectionHeader(context, '關於'),
          ListTile(
            title: const Text('版本'),
            subtitle: const Text('1.0.0'),
            trailing: const Icon(Icons.info_outline),
            onTap: () {
              _showAboutDialog(context);
            },
          ),
        ],
      ),
    );
  }
  
  // 建立區段標題
  Widget _buildSectionHeader(BuildContext context, String title) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0),
      child: Text(
        title,
        style: TextStyle(
          fontSize: 18,
          fontWeight: FontWeight.bold,
          color: Theme.of(context).colorScheme.primary,
        ),
      ),
    );
  }
  
  // 建立主題設定
  Widget _buildThemeSetting(BuildContext context, WidgetRef ref, String currentTheme) {
    return ListTile(
      title: const Text('主題'),
      subtitle: Text(_getThemeText(currentTheme)),
      trailing: DropdownButton<String>(
        value: currentTheme,
        onChanged: (value) {
          if (value != null) {
            ref.read(settingsProvider.notifier).updateSetting('theme', value);
          }
        },
        items: const [
          DropdownMenuItem(
            value: 'system',
            child: Text('跟隨系統'),
          ),
          DropdownMenuItem(
            value: 'light',
            child: Text('淺色'),
          ),
          DropdownMenuItem(
            value: 'dark',
            child: Text('深色'),
          ),
        ],
      ),
    );
  }
  
  // 建立 API URL 設定
  Widget _buildApiUrlSetting(BuildContext context, WidgetRef ref, String currentUrl) {
    return ListTile(
      title: const Text('API URL'),
      subtitle: Text(currentUrl),
      trailing: IconButton(
        icon: const Icon(Icons.edit),
        onPressed: () {
          _showApiUrlDialog(context, ref, currentUrl);
        },
      ),
    );
  }
  
  // 建立日誌保留設定
  Widget _buildLogRetentionSetting(BuildContext context, WidgetRef ref, int days) {
    return ListTile(
      title: const Text('日誌保留天數'),
      subtitle: Text('$days 天'),
      trailing: DropdownButton<int>(
        value: days,
        onChanged: (value) {
          if (value != null) {
            ref.read(settingsProvider.notifier).updateSetting('logRetention', value);
          }
        },
        items: const [
          DropdownMenuItem(
            value: 7,
            child: Text('7 天'),
          ),
          DropdownMenuItem(
            value: 14,
            child: Text('14 天'),
          ),
          DropdownMenuItem(
            value: 30,
            child: Text('30 天'),
          ),
          DropdownMenuItem(
            value: 90,
            child: Text('90 天'),
          ),
        ],
      ),
    );
  }
  
  // 建立顯示除錯日誌設定
  Widget _buildShowDebugLogsSetting(BuildContext context, WidgetRef ref, bool show) {
    return SwitchListTile(
      title: const Text('顯示除錯日誌'),
      subtitle: const Text('在日誌列表中顯示 DEBUG 級別的日誌'),
      value: show,
      onChanged: (value) {
        ref.read(settingsProvider.notifier).updateSetting('showDebugLogs', value);
      },
    );
  }
  
  // 建立啟用通知設定
  Widget _buildEnableNotificationsSetting(BuildContext context, WidgetRef ref, bool enable) {
    return SwitchListTile(
      title: const Text('啟用通知'),
      subtitle: const Text('接收錯誤和警告的通知'),
      value: enable,
      onChanged: (value) {
        ref.read(settingsProvider.notifier).updateSetting('enableNotifications', value);
      },
    );
  }
  
  // 顯示 API URL 對話框
  void _showApiUrlDialog(BuildContext context, WidgetRef ref, String currentUrl) {
    final controller = TextEditingController(text: currentUrl);
    
    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: const Text('設定 API URL'),
          content: TextField(
            controller: controller,
            decoration: const InputDecoration(
              labelText: 'API URL',
              hintText: 'http://localhost:8080/api/v1',
            ),
          ),
          actions: [
            TextButton(
              onPressed: () {
                Navigator.pop(context);
              },
              child: const Text('取消'),
            ),
            TextButton(
              onPressed: () {
                final newUrl = controller.text.trim();
                if (newUrl.isNotEmpty) {
                  ref.read(settingsProvider.notifier).updateSetting('apiUrl', newUrl);
                }
                Navigator.pop(context);
              },
              child: const Text('儲存'),
            ),
          ],
        );
      },
    );
  }
  
  // 顯示關於對話框
  void _showAboutDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: const Text('關於 DevLog AI'),
          content: const Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text('DevLog AI 是一個 AI 驅動的開發者日誌與事件分析平台。'),
              SizedBox(height: 16),
              Text('版本: 1.0.0'),
              Text('© 2023 DevLog AI Team'),
            ],
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
  
  // 取得主題文字
  String _getThemeText(String theme) {
    switch (theme) {
      case 'light':
        return '淺色';
      case 'dark':
        return '深色';
      case 'system':
        return '跟隨系統';
      default:
        return '跟隨系統';
    }
  }
}
