resource "censusworkspace_destination" "test" {
  type = var.destination_type

  name = var.destination_name

  credentials = jsonencode(var.destination_credentials)
}
