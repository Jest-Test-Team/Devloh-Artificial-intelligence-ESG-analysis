import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'screens/dashboard/dashboard_screen.dart';

void main() {
  runApp(const ProviderScope(child: DevLogAIApp()));
}

class DevLogAIApp extends StatelessWidget {
  const DevLogAIApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'DevLog AI',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.blue),
        useMaterial3: true,
        visualDensity: VisualDensity.adaptivePlatformDensity,
      ),
      home: const DashboardScreen(),
    );
  }
}
