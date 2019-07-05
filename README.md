# Snowplow Terraform Provider

[![Build Status][travis-image]][travis] [![Go Report Card][goreport-image]][goreport] [![Release][release-image]][releases] [![License][license-image]][license]

## Overview

Terraform provider for emitting Snowplow events.

## Quick start

Assuming git, **[Vagrant][vagrant-install]** and **[VirtualBox][virtualbox-install]** installed:

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

### Publishing

This is handled through CI/CD on Travis. However all binaries will be generated by using the `make` command for local publishing.

### Copyright and license

Snowplow Terraform Provider is copyright 2019 Snowplow Analytics Ltd.

Licensed under the **[Apache License, Version 2.0][license]** (the "License");
you may not use this software except in compliance with the License.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[travis-image]: https://travis-ci.com/snowplow-devops/terraform-provider-snowplow.png?branch=master
[travis]: https://travis-ci.com/snowplow-devops/terraform-provider-snowplow

[release-image]: http://img.shields.io/badge/release-0.1.0-6ad7e5.svg?style=flat
[releases]: https://github.com/snowplow-devops/terraform-provider-snowplow/releases

[license-image]: http://img.shields.io/badge/license-Apache--2-blue.svg?style=flat
[license]: http://www.apache.org/licenses/LICENSE-2.0

[goreport-image]: https://goreportcard.com/badge/github.com/snowplow-devops/terraform-provider-snowplow
[goreport]: https://goreportcard.com/report/github.com/snowplow-devops/terraform-provider-snowplow

[vagrant-url]: http://docs.vagrantup.com/v2/installation/index.html
[virtualbox-url]: https://www.virtualbox.org/wiki/Downloads