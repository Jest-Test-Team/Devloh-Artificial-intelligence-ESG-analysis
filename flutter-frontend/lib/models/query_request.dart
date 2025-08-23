import 'package:json_annotation/json_annotation.dart';

part 'query_request.g.dart';

@JsonSerializable()
class QueryRequest {
  final String query;
  @JsonKey(name: 'project_id')
  final String? projectId;
  final Map<String, dynamic>? options;

  QueryRequest({
    required this.query,
    this.projectId,
    this.options,
  });

  factory QueryRequest.fromJson(Map<String, dynamic> json) => _$QueryRequestFromJson(json);
  Map<String, dynamic> toJson() => _$QueryRequestToJson(this);
  
  // 創建一個帶有預設選項的查詢請求
  factory QueryRequest.withDefaultOptions({
    required String query,
    String? projectId,
  }) {
    return QueryRequest(
      query: query,
      projectId: projectId,
      options: {
        'temperature': 0.0,
        'max_tokens': 1000,
        'include_sources': true,
      },
    );
  }
}