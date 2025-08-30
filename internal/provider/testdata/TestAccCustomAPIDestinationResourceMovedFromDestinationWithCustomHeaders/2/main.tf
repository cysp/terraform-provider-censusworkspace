moved {
  from = censusworkspace_destination.test
  to   = censusworkspace_custom_api_destination.test
}

resource "censusworkspace_custom_api_destination" "test" {
  name = var.destination_name

  credentials = var.destination_credentials
}
