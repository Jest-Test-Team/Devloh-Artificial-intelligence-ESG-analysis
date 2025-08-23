provider "oci" {
  # Authentication details will be provided by environment variables or ~/.oci/config
  region = var.region
}

module "networking" {
  source = "../modules/networking"
  
  compartment_id = var.compartment_id
  vcn_cidr       = var.vcn_cidr
  subnet_cidr    = var.subnet_cidr
}

module "oke-cluster" {
  source = "../modules/oke-cluster"
  
  compartment_id        = var.compartment_id
  kubernetes_version    = var.kubernetes_version
  vcn_id                = module.networking.vcn_id
  subnet_id             = module.networking.subnet_id
  node_pool_shape       = var.node_pool_shape
  node_pool_image_id    = var.node_pool_image_id
  node_pool_size        = var.node_pool_size
}

module "database" {
  source = "../modules/database"
  
  compartment_id = var.compartment_id
  db_name        = var.db_name
  db_password    = var.db_password
  subnet_id      = module.networking.subnet_id
}
