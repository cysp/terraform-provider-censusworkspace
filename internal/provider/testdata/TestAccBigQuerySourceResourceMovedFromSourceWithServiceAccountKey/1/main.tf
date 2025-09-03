resource "censusworkspace_source" "test" {
  type  = "big_query"
  label = var.source_label

  credentials = jsonencode(var.source_credentials)
}
