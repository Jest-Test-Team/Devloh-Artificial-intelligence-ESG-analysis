variable "compartment_id" {
  description = "OCI compartment OCID"
  type        = string
}

variable "kubernetes_version" {
  description = "Kubernetes version for OKE cluster"
  type        = string
}

variable "vcn_id" {
  description = "ID of the VCN"
  type        = string
}

variable "subnet_id" {
  description = "ID of the subnet"
  type        = string
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
