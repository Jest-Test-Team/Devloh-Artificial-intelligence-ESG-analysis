output "vcn_id" {
  description = "ID of the created VCN"
  value       = oci_core_vcn.devlog_ai_vcn.id
}

output "subnet_id" {
  description = "ID of the created subnet"
  value       = oci_core_subnet.public_subnet.id
}

output "internet_gateway_id" {
  description = "ID of the created internet gateway"
  value       = oci_core_internet_gateway.internet_gateway.id
}
