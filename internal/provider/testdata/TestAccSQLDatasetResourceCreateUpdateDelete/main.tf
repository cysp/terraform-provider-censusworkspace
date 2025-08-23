resource "censusworkspace_sql_dataset" "test" {
  name = var.dataset_name

  source_id = var.dataset_source_id

  query = var.dataset_query
}
