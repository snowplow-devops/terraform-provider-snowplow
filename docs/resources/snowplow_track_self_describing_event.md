# TrackSelfDescribingEvent Resource

Emits an event to the configured collector upon creation, update, or deletion of the resource.

## Example Usage

```hcl
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

## Argument Reference

* `create_event` - Event emmitted during creation of this TF resource
* `update_event` - Event emmitted during update of this TF resource
* `delete_event` - Event emmitted during deletion of this TF resource
* `contexts` - A payload containing additional context