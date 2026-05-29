resource "censusworkspace_big_query_source" "test" {
  name = var.source_name

  credentials = var.source_credentials

  warehouse_writeback_retention_in_days = var.source_warehouse_writeback_retention_in_days
}
