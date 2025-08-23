variable "region" {
  description = "OCI region"
  type        = string
  default     = "us-phoenix-1"
}

variable "compartment_id" {
  description = "OCI compartment OCID"
  type        = string
}

variable "vcn_cidr" {
  description = "CIDR block for the VCN"
  type        = string
  default     = "10.0.0.0/16"
}

variable "subnet_cidr" {
  description = "CIDR block for the subnet"
  type        = string
  default     = "10.0.1.0/24"
}

variable "kubernetes_version" {
  description = "Kubernetes version for OKE cluster"
  type        = string
  default     = "v1.27.0"
}

variable "node_pool_shape" {
  description = "Shape for the node pool instances"
  type        = string
  default     = "VM.Standard.A1.Flex"
}

variable "node_pool_image_id" {
  description = "Image ID for the node pool instances"
  type        = string
}

variable "node_pool_size" {
  description = "Number of nodes in the node pool"
  type        = number
  default     = 2
}

variable "db_name" {
  description = "Name for the database"
  type        = string
  default     = "devlogaidb"
}

variable "db_password" {
  description = "Password for the database"
  type        = string
  sensitive   = true
}
