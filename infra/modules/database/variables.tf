variable "compartment_id" {
  description = "OCI compartment OCID"
  type        = string
}

variable "db_name" {
  description = "Name for the database"
  type        = string
}

variable "db_password" {
  description = "Password for the database"
  type        = string
  sensitive   = true
}

variable "subnet_id" {
  description = "ID of the subnet"
  type        = string
}
