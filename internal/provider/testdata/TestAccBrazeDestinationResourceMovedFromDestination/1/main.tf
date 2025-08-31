resource "censusworkspace_destination" "test" {
  type = "braze"

  name = var.destination_name

  credentials = jsonencode(var.destination_credentials)
}
