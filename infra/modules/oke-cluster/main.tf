resource "oci_containerengine_cluster" "oke_cluster" {
  compartment_id     = var.compartment_id
  kubernetes_version = var.kubernetes_version
  name               = "DevLogAI-OKE-Cluster"
  vcn_id             = var.vcn_id
  
  options {
    service_lb_subnet_ids = [var.subnet_id]
    
    add_ons {
      is_kubernetes_dashboard_enabled = true
      is_tiller_enabled               = false
    }
    
    kubernetes_network_config {
      pods_cidr     = "10.244.0.0/16"
      services_cidr = "10.96.0.0/16"
    }
  }
}

resource "oci_containerengine_node_pool" "oke_node_pool" {
  cluster_id         = oci_containerengine_cluster.oke_cluster.id
  compartment_id     = var.compartment_id
  kubernetes_version = var.kubernetes_version
  name               = "DevLogAI-Node-Pool"
  node_shape         = var.node_pool_shape
  
  node_config_details {
    size = var.node_pool_size
    
    placement_configs {
      availability_domain = data.oci_identity_availability_domains.ads.availability_domains[0].name
      subnet_id           = var.subnet_id
    }
  }
  
  node_source_details {
    image_id    = var.node_pool_image_id
    source_type = "IMAGE"
  }
  
  dynamic "node_shape_config" {
    for_each = var.node_pool_shape == "VM.Standard.A1.Flex" ? [1] : []
    content {
      memory_in_gbs = 24
      ocpus         = 4
    }
  }
}

// Data source to get availability domains
data "oci_identity_availability_domains" "ads" {
  compartment_id = var.compartment_id
}
