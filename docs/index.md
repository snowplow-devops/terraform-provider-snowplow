# Snowplow Provider

Terraform provider for emitting Snowplow events

## Example Usage ##

To actually start tracking Snowplow events from Terraform you will need to configure the `provider` and a `resource`:

```hcl
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
```

Now that the provider is configured we can track an event!

#### How to track: `self_describing_event` ####

```hcl
locals {
  tf_module_context = {
    name    = "aws_module"
    version = "1.1.2"
  }
}

resource "snowplow_track_self_describing_event" "module_action" {
  create_event = {
    iglu_uri = "iglu:com.acme/lifecycle/jsonschema/1-0-0",
    payload  = "{\"actionType\":\"create\"}"
  }

  update_event = {
    iglu_uri = "iglu:com.acme/lifecycle/jsonschema/1-0-0",
    payload  = "{\"actionType\":\"update\"}"
  }

  delete_event = {
    iglu_uri = "iglu:com.acme/lifecycle/jsonschema/1-0-0",
    payload  = "{\"actionType\":\"delete\"}"
  }

  contexts = [
    {
      iglu_uri = "iglu:com.acme/module_context/jsonschema/1-0-0",
      payload  = "${jsonencode(local.tf_module_context)}",
    },
  ]
}
```

## Argument Reference ##

* `collector_uri` (Required) URI of your Snowplow Collector
* `tracker_app_id` (Optional) Optional application ID (Default: "")
* `tracker_namespace` (Optional) Optional namespace (Default: "")
* `tracker_platform` (Optional) Optional platform (Default: srv)
* `emitter_request_type` (Optional) Whether to use GET or POST requests to emit events (Default: "POST")
* `emitter_protocol` (Optional) Whether to use HTTP or HTTPS to send events (Default: "HTTPS")