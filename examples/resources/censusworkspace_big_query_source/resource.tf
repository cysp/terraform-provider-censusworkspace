resource "censusworkspace_big_query_source" "test" {
  sync_engine = "advanced"

  name = "bigquery_project_id"

  credentials = {
    project_id = "project-id"
    location   = "US"
  }
}
