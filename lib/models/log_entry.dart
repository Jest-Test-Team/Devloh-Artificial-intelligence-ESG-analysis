import 'package:json_annotation/json_annotation.dart';

part 'log_entry.g.dart';

@JsonSerializable()
class LogEntry {
  final String? id;
  final String timestamp;
  final String level;
  final String message;
  final String source;
  final String serviceName;
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

  factory LogEntry.fromJson(Map<String, dynamic> json) => _$LogEntryFromJson(json);
  Map<String, dynamic> toJson() => _$LogEntryToJson(this);
}
