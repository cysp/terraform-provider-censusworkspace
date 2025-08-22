resource "censusworkspace_destination" "test" {
  name = "custom"

  type = "custom_api"

  credentials = jsonencode({
    webhook_url = "https://example.org/census-destination"
  })
}
