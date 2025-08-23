resource "oci_core_vcn" "devlog_ai_vcn" {
  compartment_id = var.compartment_id
  cidr_blocks    = [var.vcn_cidr]
  display_name   = "DevLogAI-VCN"
  dns_label      = "devlogaivcn"
}

resource "oci_core_internet_gateway" "internet_gateway" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.devlog_ai_vcn.id
  display_name   = "DevLogAI-InternetGateway"
  enabled        = true
}

resource "oci_core_route_table" "public_route_table" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.devlog_ai_vcn.id
  display_name   = "DevLogAI-PublicRouteTable"
  
  route_rules {
    destination       = "0.0.0.0/0"
    network_entity_id = oci_core_internet_gateway.internet_gateway.id
  }
}

resource "oci_core_subnet" "public_subnet" {
  cidr_block        = var.subnet_cidr
  compartment_id    = var.compartment_id
  vcn_id            = oci_core_vcn.devlog_ai_vcn.id
  display_name      = "DevLogAI-PublicSubnet"
  route_table_id    = oci_core_route_table.public_route_table.id
  security_list_ids = [oci_core_security_list.security_list.id]
  dns_label         = "devlogainet"
}

resource "oci_core_security_list" "security_list" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.devlog_ai_vcn.id
  display_name   = "DevLogAI-SecurityList"
  
  // Allow all outbound traffic
  egress_security_rules {
    destination = "0.0.0.0/0"
    protocol    = "all"
    stateless   = false
  }
  
  // Allow SSH inbound
  ingress_security_rules {
    protocol  = "6" // TCP
    source    = "0.0.0.0/0"
    stateless = false
    
    tcp_options {
      min = 22
      max = 22
    }
  }
  
  // Allow HTTP inbound
  ingress_security_rules {
    protocol  = "6" // TCP
    source    = "0.0.0.0/0"
    stateless = false
    
    tcp_options {
      min = 80
      max = 80
    }
  }
  
  // Allow HTTPS inbound
  ingress_security_rules {
    protocol  = "6" // TCP
    source    = "0.0.0.0/0"
    stateless = false
    
    tcp_options {
      min = 443
      max = 443
    }
  }
  
  // Allow Kubernetes API server
  ingress_security_rules {
    protocol  = "6" // TCP
    source    = "0.0.0.0/0"
    stateless = false
    
    tcp_options {
      min = 6443
      max = 6443
    }
  }
}
