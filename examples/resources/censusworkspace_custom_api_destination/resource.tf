resource "censusworkspace_custom_api_destination" "this" {
  name = "destination name"

  credentials = {
    api_version = 1
    webhook_url = "https://example.org/census-destination"
    custom_headers = {
      "x-client-id" = {
        value = "123"
      }
      "x-client-secret" = {
        value     = "secret"
        is_secret = true
      }
    }
  }
}
