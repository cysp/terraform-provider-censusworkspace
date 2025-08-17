resource "censusworkspace_source" "test" {
  type = var.source_type

  label = var.source_label

  credentials = jsonencode(var.source_credentials)
}
