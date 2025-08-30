resource "censusworkspace_custom_api_destination" "test" {
  name = var.destination_name

  credentials = var.destination_credentials
}
