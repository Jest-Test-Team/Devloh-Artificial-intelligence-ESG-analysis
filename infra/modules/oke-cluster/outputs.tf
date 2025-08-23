output "cluster_id" {
  description = "ID of the created OKE cluster"
  value       = oci_containerengine_cluster.oke_cluster.id
}

output "cluster_endpoint" {
  description = "Endpoint of the Kubernetes API server"
  value       = oci_containerengine_cluster.oke_cluster.endpoints[0].public_endpoint
}

output "node_pool_id" {
  description = "ID of the created node pool"
  value       = oci_containerengine_node_pool.oke_node_pool.id
}
