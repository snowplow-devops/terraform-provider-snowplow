//
// Copyright (c) 2019-2023 Snowplow Analytics Ltd. All rights reserved.
//
// This program is licensed to you under the Apache License Version 2.0,
// and you may not use this file except in compliance with the Apache License Version 2.0.
// You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the Apache License Version 2.0 is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
//

package main

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/assert"
)

func TestInitTracker(t *testing.T) {
	assert := assert.New(t)

	// Setup Tracker
	ctx := SnowplowProviderModel{
		CollectorURI:       types.StringValue("com.acme"),
		TrackerAppID:       types.StringValue(""),
		TrackerNamespace:   types.StringValue(""),
		TrackerPlatform:    types.StringValue("srv"),
		EmitterRequestType: types.StringValue("GET"),
		EmitterProtocol:    types.StringValue("HTTP"),
	}
	ctxR := TrackSelfDescribingEventResourceModel{
		CollectorURI:       types.StringValue(""),
		TrackerAppID:       types.StringValue(""),
		TrackerNamespace:   types.StringValue(""),
		TrackerPlatform:    types.StringValue(""),
		EmitterRequestType: types.StringValue(""),
		EmitterProtocol:    types.StringValue(""),
	}

	trackerChan := make(chan int, 1)
	tracker, err := InitTracker(ctx, ctxR, trackerChan)
	assert.NotNil(tracker)
	assert.Nil(err)
	assert.Equal("http://com.acme/i", tracker.Emitter.GetCollectorUrl())
}

func TestInitTracker_WithOverrides(t *testing.T) {
	assert := assert.New(t)

	// Setup Tracker
	ctx := SnowplowProviderModel{
		CollectorURI:       types.StringValue("com.acme"),
		TrackerAppID:       types.StringValue(""),
		TrackerNamespace:   types.StringValue(""),
		TrackerPlatform:    types.StringValue("srv"),
		EmitterRequestType: types.StringValue("GET"),
		EmitterProtocol:    types.StringValue("HTTP"),
	}
	ctxR := TrackSelfDescribingEventResourceModel{
		CollectorURI:       types.StringValue("com.acme.override"),
		TrackerAppID:       types.StringValue(""),
		TrackerNamespace:   types.StringValue(""),
		TrackerPlatform:    types.StringValue(""),
		EmitterRequestType: types.StringValue(""),
		EmitterProtocol:    types.StringValue(""),
	}

	trackerChan := make(chan int, 1)
	tracker, err := InitTracker(ctx, ctxR, trackerChan)
	assert.NotNil(tracker)
	assert.Nil(err)
	assert.Equal("http://com.acme.override/i", tracker.Emitter.GetCollectorUrl())
}

func TestInitTracker_WithEmptyCollectorURI(t *testing.T) {
	assert := assert.New(t)

	// Setup Tracker
	ctx := SnowplowProviderModel{
		CollectorURI:       types.StringValue(""),
		TrackerAppID:       types.StringValue(""),
		TrackerNamespace:   types.StringValue(""),
		TrackerPlatform:    types.StringValue("srv"),
		EmitterRequestType: types.StringValue("GET"),
		EmitterProtocol:    types.StringValue("HTTP"),
	}
	ctxR := TrackSelfDescribingEventResourceModel{
		CollectorURI:       types.StringValue(""),
		TrackerAppID:       types.StringValue(""),
		TrackerNamespace:   types.StringValue(""),
		TrackerPlatform:    types.StringValue(""),
		EmitterRequestType: types.StringValue(""),
		EmitterProtocol:    types.StringValue(""),
	}

	trackerChan := make(chan int, 1)
	tracker, err := InitTracker(ctx, ctxR, trackerChan)
	assert.Nil(tracker)
	assert.NotNil(err)
	assert.Equal("URI of the Snowplow Collector is empty - this can be set either at the provider or resource level with the 'collector_uri' input", err.Error())
}

func TestResource_UpgradeFromVersion(t *testing.T) {
	config := `resource "snowplow_track_self_describing_event" "example" {
		collector_uri = "localhost:9090"
		emitter_protocol = "HTTP"
		contexts = []
		create_event = {
			iglu_uri = "iglu:com.example/create/jsonschema/1-0-0"
			payload = jsonencode({
				empty = true
			})
		}
		update_event = {
			iglu_uri = "iglu:com.example/update/jsonschema/1-0-0"
			payload = jsonencode({
				empty = true
			})
		}
		delete_event = {
			iglu_uri = "iglu:com.example/delete/jsonschema/1-0-0"
			payload = jsonencode({
				empty = true
			})
		}
	}`
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowplow": {
						VersionConstraint: "0.7.3",
						Source:            "snowplow-devops/snowplow",
					},
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowplow_track_self_describing_event.example", "collector_uri", "localhost:9090"),
					resource.TestCheckResourceAttr("snowplow_track_self_describing_event.example", "emitter_protocol", "HTTP"),
				),
			},
			{
				ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
					"snowplow": providerserver.NewProtocol5WithError(NewProvider("dev")()),
				},
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"snowplow_track_self_describing_event.example",
							plancheck.ResourceActionNoop,
						),
					},
				},
			},
		},
	})
}
