resource "censusworkspace_destination" "test" {
  type = "custom_api"

  name = var.destination_name

  credentials = jsonencode(var.destination_credentials)
}
