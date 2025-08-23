import 'package:json_annotation/json_annotation.dart';

part 'query_response.g.dart';

@JsonSerializable()
class QueryResponse {
  final String answer;
  final List<String>? relatedLogs;
  final List<String>? sources;
  final int processingMs;

  QueryResponse({
    required this.answer,
    this.relatedLogs,
    this.sources,
    required this.processingMs,
  });

  factory QueryResponse.fromJson(Map<String, dynamic> json) => _$QueryResponseFromJson(json);
  Map<String, dynamic> toJson() => _$QueryResponseToJson(this);
}
