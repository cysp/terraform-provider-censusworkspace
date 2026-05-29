resource "censusworkspace_big_query_source" "test" {
  name = var.source_name

  credentials = var.source_credentials
}
