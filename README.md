# Snowplow Terraform Provider

[![Build Status][travis-image]][travis] [![Go Report Card][goreport-image]][goreport] [![Release][release-image]][releases] [![License][license-image]][license]

## Overview

Terraform provider for emitting Snowplow events.

## Quick start

Assuming git, **[Vagrant][vagrant-url]** and **[VirtualBox][virtualbox-url]** installed:

```bash
 host> git clone https://github.com/snowplow-devops/terraform-provider-snowplow
 host> cd terraform-provider-snowplow
 host> vagrant up && vagrant ssh
guest> cd /opt/gopath/src/github.com/snowplow/terraform-provider-snowplow
guest> make test
guest> make
```

To remove all build files:

```bash
guest> make clean
```

To format the golang code in the source directory:

```bash
guest> make format
```

**Note:** Always run `format` before submitting any code.

**Note:** The `make test` command also generates a code coverage file which can be found at `build/coverage/coverage.html`.

## Installation

First download the pre-compiled binary for your platform from our Bintray at the following links or generate the binaries locally using the provided `make` command:

* [Darwin (macOS)](https://bintray.com/snowplow/snowplow-generic/download_file?file_path=terraform_provider_snowplow_0.1.1_darwin_amd64.zip)
* [Linux](https://bintray.com/snowplow/snowplow-generic/download_file?file_path=terraform_provider_snowplow_0.1.1_linux_amd64.zip)
* [Windows](https://bintray.com/snowplow/snowplow-generic/download_file?file_path=terraform_provider_snowplow_0.1.1_windows_amd64.zip)

Once downloaded "unzip" to extract the binary which should be called `terraform-provider-snowplow_v0.1.1`.

From here you will need to move the binary into your Terraform plugins directory - depending on your platform / installation this might change but generally speaking they are located at:

* Darwin & Linux: `~/.terraform.d/plugins`
* Windows: `%APPDATA%\terraform.d\plugins`

## How to use?

### Setting up the provider

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

### Tracking events

Currently only one tracking resource has been exposed in the form of a `SelfDescribing Event` - the other pre-built options do not really make sense in a Terraform world!

#### How to track: `self_describing_event`

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

To get around nested dictionary limitations in Terraform `0.11.x` we are using stringified JSON as inputs for the payloads.  In the example above we are leveraging the `jsonencode` builtin function to handle converting a Terraform dictionary into a JSON but you can also hand-craft encoded JSON strings if you prefer.

The above function lets you define a single event which is coupled with as many contexts as you would like to attach - when you run `terraform apply` it will send an event to your defined collector and on a `200 OK` response code from the collector log this as a successful resource creation.

Its important to define a base event for _each_ part of the resource lifecycle so that you can differentiate your events later and to be able to reason about where in the lifecycle a resource might be.

### How to stop Terraform turning primitives into strings

We are using a lot of stringified JSON as input above - the `jsonencode` function unfortunately turns a lot of our primitives ("floats", "ints" and "booleans") into strings through this conversion.  One way to get around this is to use regex to "fix" the payload after conversion:

```hcl
# Convert int, floats
payload_1 = "${replace(jsonencode(var.payload), "/\"([0-9]+\\.?[0-9]*)\"/", "$1")}"

# Convert "true" > true
payload_2 = "${replace(local.payload_1, "/\"(true)\"/", "$1")}"

# Convert "false" > false
payload_3 = "${replace(local.payload_2, "/\"(false)\"/", "$1")}"
```
 
### Publishing

This is handled through CI/CD on Travis. However all binaries will be generated by using the `make` command for local publishing.

### Copyright and license

Snowplow Terraform Provider is copyright 2019-2020 Snowplow Analytics Ltd.

Licensed under the **[Apache License, Version 2.0][license]** (the "License");
you may not use this software except in compliance with the License.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[travis-image]: https://travis-ci.com/snowplow-devops/terraform-provider-snowplow.png?branch=master
[travis]: https://travis-ci.com/snowplow-devops/terraform-provider-snowplow

[release-image]: http://img.shields.io/badge/release-0.1.1-6ad7e5.svg?style=flat
[releases]: https://github.com/snowplow-devops/terraform-provider-snowplow/releases

[license-image]: http://img.shields.io/badge/license-Apache--2-blue.svg?style=flat
[license]: http://www.apache.org/licenses/LICENSE-2.0

[goreport-image]: https://goreportcard.com/badge/github.com/snowplow-devops/terraform-provider-snowplow
[goreport]: https://goreportcard.com/report/github.com/snowplow-devops/terraform-provider-snowplow

[vagrant-url]: http://docs.vagrantup.com/v2/installation/index.html
[virtualbox-url]: https://www.virtualbox.org/wiki/Downloads
