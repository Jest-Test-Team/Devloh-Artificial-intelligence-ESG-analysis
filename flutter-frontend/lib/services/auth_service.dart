import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:shared_preferences/shared_preferences.dart';

// 提供 AuthService 的 Provider
final authServiceProvider = Provider<AuthService>((ref) {
  return AuthService();
});

// 提供身份驗證狀態的 Provider
final authStateProvider = StateNotifierProvider<AuthStateNotifier, bool>((ref) {
  final authService = ref.watch(authServiceProvider);
  return AuthStateNotifier(authService);
});

// 身份驗證狀態管理
class AuthStateNotifier extends StateNotifier<bool> {
  final AuthService _authService;

  AuthStateNotifier(this._authService) : super(false) {
    // 初始化時檢查登入狀態
    _checkAuthState();
  }

  Future<void> _checkAuthState() async {
    final isLoggedIn = await _authService.isLoggedIn();
    state = isLoggedIn;
  }

  Future<bool> login(String email, String password) async {
    final result = await _authService.login(email, password);
    if (result) {
      state = true;
    }
    return result;
  }

  Future<void> logout() async {
    await _authService.logout();
    state = false;
  }
}

// 身份驗證服務
class AuthService {
  final FlutterSecureStorage _secureStorage = const FlutterSecureStorage();
  static const String _tokenKey = 'auth_token';
  static const String _userIdKey = 'user_id';
  static const String _userNameKey = 'user_name';
  static const String _userEmailKey = 'user_email';
  static const String _projectIdKey = 'project_id';

  // 初始化服務
  Future<void> initialize() async {
    // 可以在這裡進行額外的初始化工作
  }

  // 檢查是否已登入
  Future<bool> isLoggedIn() async {
    final token = await _secureStorage.read(key: _tokenKey);
    return token != null;
  }

  // 登入
  Future<bool> login(String email, String password) async {
    try {
      // 這裡應該是實際的 API 呼叫
      // 目前使用模擬資料
      if (email == 'test@example.com' && password == 'password') {
        await _secureStorage.write(key: _tokenKey, value: 'mock_token_123');
        await _secureStorage.write(key: _userIdKey, value: 'user_1');
        await _secureStorage.write(key: _userNameKey, value: 'Test User');
        await _secureStorage.write(key: _userEmailKey, value: email);
        await _secureStorage.write(key: _projectIdKey, value: 'project_1');
        
        // 儲存最後登入時間
        final prefs = await SharedPreferences.getInstance();
        await prefs.setString('last_login', DateTime.now().toIso8601String());
        
        return true;
      }
      return false;
    } catch (e) {
      return false;
    }
  }

  // 登出
  Future<void> logout() async {
    await _secureStorage.delete(key: _tokenKey);
    await _secureStorage.delete(key: _userIdKey);
    await _secureStorage.delete(key: _userNameKey);
    await _secureStorage.delete(key: _userEmailKey);
    await _secureStorage.delete(key: _projectIdKey);
  }

  // 取得使用者資料
  Future<Map<String, String?>> getUserData() async {
    final userId = await _secureStorage.read(key: _userIdKey);
    final userName = await _secureStorage.read(key: _userNameKey);
    final userEmail = await _secureStorage.read(key: _userEmailKey);
    final projectId = await _secureStorage.read(key: _projectIdKey);

    return {
      'userId': userId,
      'userName': userName,
      'userEmail': userEmail,
      'projectId': projectId,
    };
  }

  // 取得授權 Token
  Future<String?> getAuthToken() async {
    return _secureStorage.read(key: _tokenKey);
  }
}
