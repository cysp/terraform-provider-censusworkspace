resource "censusworkspace_source" "test" {
  type = var.source_type

  name = var.source_name

  credentials = jsonencode(var.source_credentials)

  warehouse_writeback_retention_in_days = var.source_warehouse_writeback_retention_in_days
}
