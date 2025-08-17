resource "censusworkspace_source" "test" {
  type = "big_query"

  label = "BigQuery - project-id"

  credentials = jsonencode({
    project_id = "project-id"
    location   = "US"
  })
}
