resource "censusworkspace_source" "test" {
  type = var.source_type

  label = var.source_label

  credentials = jsonencode(var.source_credentials)

  warehouse_writeback_retention_in_days = var.source_warehouse_writeback_retention_in_days
}
