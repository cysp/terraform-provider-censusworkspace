variable "source_type" {
  type = string
}

variable "source_label" {
  type    = string
  default = null
}

variable "source_credentials" {
}

variable "source_warehouse_writeback_retention_in_days" {
  type    = number
  default = null
}
