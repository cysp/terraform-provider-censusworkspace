resource "censusworkspace_source" "test" {
  type  = "big_query"
  name = var.source_name

  credentials = jsonencode(var.source_credentials)
}
