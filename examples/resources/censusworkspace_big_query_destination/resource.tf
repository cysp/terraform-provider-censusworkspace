resource "censusworkspace_big_query_destination" "test" {
  name = "BigQuery - project-id"

  credentials = {
    project_id = "project-id"
    location   = "US"
  }
}
