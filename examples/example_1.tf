# Minimal configuration
provider "snowplow" {
  collector_uri = "com.acme.collector"
}

# Optional extra configuration options
provider "snowplow" {
  collector_uri        = "com.acme.collector"
  tracker_app_id       = "terraform" # Default; ""
  tracker_namespace    = "terraform" # Default; ""
  tracker_platform     = "mob" # Default; srv (server)
  emitter_request_type = "GET" # Default; POST
  emitter_protocol     = "HTTP" # Default; HTTPS
}
