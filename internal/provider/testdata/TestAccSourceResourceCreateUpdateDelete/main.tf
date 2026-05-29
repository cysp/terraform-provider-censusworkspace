resource "censusworkspace_source" "test" {
  type = var.source_type

  name = var.source_name

  credentials = jsonencode(var.source_credentials)
}
