moved {
  from = censusworkspace_source.test
  to   = censusworkspace_big_query_source.test
}

resource "censusworkspace_big_query_source" "test" {
  name = var.source_name

  credentials = var.source_credentials
}
