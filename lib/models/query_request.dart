import 'package:json_annotation/json_annotation.dart';

part 'query_request.g.dart';

@JsonSerializable()
class QueryRequest {
  final String query;
  final String? projectId;

  QueryRequest({
    required this.query,
    this.projectId,
  });

  factory QueryRequest.fromJson(Map<String, dynamic> json) => _$QueryRequestFromJson(json);
  Map<String, dynamic> toJson() => _$QueryRequestToJson(this);
}
