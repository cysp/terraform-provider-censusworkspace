moved {
  from = censusworkspace_source.test
  to   = censusworkspace_big_query_source.test
}

resource "censusworkspace_big_query_source" "test" {
  label = var.source_label

  credentials = var.source_credentials
}
