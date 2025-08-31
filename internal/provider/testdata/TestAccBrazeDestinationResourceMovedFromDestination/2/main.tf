moved {
  from = censusworkspace_destination.test
  to   = censusworkspace_braze_destination.test
}

resource "censusworkspace_braze_destination" "test" {
  name = var.destination_name

  credentials = var.destination_credentials
}
