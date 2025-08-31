resource "censusworkspace_destination" "test" {
  type = "big_query"

  name = var.destination_name

  credentials = jsonencode(var.destination_credentials)
}
