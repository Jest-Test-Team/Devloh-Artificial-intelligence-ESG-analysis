import 'package:json_annotation/json_annotation.dart';

part 'log_entry.g.dart';

@JsonSerializable()
class LogEntry {
  final String? id;
  final String timestamp;
  final String level;
  final String message;
  final String source;
  @JsonKey(name: 'service_name')
  final String serviceName;
  @JsonKey(name: 'project_id')
  final String projectId;
  final Map<String, String>? metadata;

  LogEntry({
    this.id,
    required this.timestamp,
    required this.level,
    required this.message,
    required this.source,
    required this.serviceName,
    required this.projectId,
    this.metadata,
  });

  factory LogEntry.fromJson(Map<String, dynamic> json) =>
      _$LogEntryFromJson(json);
  Map<String, dynamic> toJson() => _$LogEntryToJson(this);

  // 創建一個帶有當前時間戳的新日誌條目
  factory LogEntry.now({
    String? id,
    required String level,
    required String message,
    required String source,
    required String serviceName,
    required String projectId,
    Map<String, String>? metadata,
  }) {
    return LogEntry(
      id: id,
      timestamp: DateTime.now().toUtc().toIso8601String(),
      level: level,
      message: message,
      source: source,
      serviceName: serviceName,
      projectId: projectId,
      metadata: metadata,
    );
  }
}
