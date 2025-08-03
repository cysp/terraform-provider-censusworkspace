resource "censusworkspace_source" "test" {
  type = "big_query"

  label = "Test Source"

  credentials = jsonencode({
    project_id = "project-id"
    location   = "US"
  })
}
