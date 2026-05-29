resource "censusworkspace_source" "test" {
  type        = "big_query"
  sync_engine = "advanced"

  name = "BigQuery - project-id"

  credentials = jsonencode({
    project_id = "project-id"
    location   = "US"
  })
}
