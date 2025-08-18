resource "censusworkspace_big_query_source" "test" {
  sync_engine = "advanced"

  label = "BigQuery - project-id"

  credentials = {
    project_id = "project-id"
    location   = "US"
  }
}
