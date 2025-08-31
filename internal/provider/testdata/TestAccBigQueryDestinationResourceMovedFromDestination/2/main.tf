moved {
  from = censusworkspace_destination.test
  to   = censusworkspace_big_query_destination.test
}

resource "censusworkspace_big_query_destination" "test" {
  name = var.destination_name

  credentials = var.destination_credentials
}
