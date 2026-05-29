resource "censusworkspace_big_query_source" "test" {
  sync_engine = "advanced"

  name = "BigQuery - project-id"

  credentials = {
    project_id = "project-id"
    location   = "US"
  }
}
