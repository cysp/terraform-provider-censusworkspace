resource "censusworkspace_source" "test" {
  type = "big_query"

  credentials = jsonencode({
    project_id = "project-id"
    location   = "US"
  })
}
