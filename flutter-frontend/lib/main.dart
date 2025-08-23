import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:devlog_ai/screens/auth/login_screen.dart';
import 'package:devlog_ai/screens/dashboard/dashboard_screen.dart';
import 'package:devlog_ai/services/auth_service.dart';
import 'package:devlog_ai/theme/app_theme.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Initialize services
  final authService = AuthService();
  await authService.initialize();
  
  runApp(
    ProviderScope(
      child: DevLogAIApp(authService: authService),
    ),
  );
}

class DevLogAIApp extends ConsumerWidget {
  final AuthService authService;
  
  const DevLogAIApp({
    required this.authService,
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final isLoggedIn = ref.watch(authStateProvider);
    
    return MaterialApp(
      title: 'DevLog AI',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.lightTheme,
      darkTheme: AppTheme.darkTheme,
      themeMode: ThemeMode.system,
      home: isLoggedIn ? const DashboardScreen() : const LoginScreen(),
    );
  }
}