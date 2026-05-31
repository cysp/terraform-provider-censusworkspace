resource "censusworkspace_big_query_source" "test" {
  label = var.source_label

  credentials = var.source_credentials

  warehouse_writeback_retention_in_days = var.source_warehouse_writeback_retention_in_days
}
