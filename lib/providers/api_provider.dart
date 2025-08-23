import 'package:flutter_riverpod/flutter_riverpod.dart';

final apiBaseUrlProvider = Provider<String>((ref) {
  // Default to localhost for development
  return 'http://localhost:8080';
});
