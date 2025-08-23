import 'package:json_annotation/json_annotation.dart';

part 'query_response.g.dart';

@JsonSerializable()
class QueryResponse {
  final String answer;
  @JsonKey(name: 'related_logs')
  final List<String>? relatedLogs;
  final List<String>? sources;
  @JsonKey(name: 'processing_ms')
  final int processingMs;

  QueryResponse({
    required this.answer,
    this.relatedLogs,
    this.sources,
    required this.processingMs,
  });

  factory QueryResponse.fromJson(Map<String, dynamic> json) =>
      _$QueryResponseFromJson(json);
  Map<String, dynamic> toJson() => _$QueryResponseToJson(this);

  // 創建一個錯誤響應
  factory QueryResponse.error(String errorMessage) {
    return QueryResponse(
      answer: 'Error: $errorMessage',
      relatedLogs: [],
      sources: [],
      processingMs: 0,
    );
  }
}
