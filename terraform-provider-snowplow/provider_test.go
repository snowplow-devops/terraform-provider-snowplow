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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestInitTracker(t *testing.T) {
	assert := assert.New(t)

	// Setup Tracker
	ctx := Context{
		CollectorURI:       "com.acme",
		TrackerAppID:       "",
		TrackerNamespace:   "",
		TrackerPlatform:    "srv",
		EmitterRequestType: "GET",
		EmitterProtocol:    "HTTP",
	}
	ctxR := Context{
		CollectorURI:       "",
		TrackerAppID:       "",
		TrackerNamespace:   "",
		TrackerPlatform:    "",
		EmitterRequestType: "",
		EmitterProtocol:    "",
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
	ctx := Context{
		CollectorURI:       "com.acme",
		TrackerAppID:       "",
		TrackerNamespace:   "",
		TrackerPlatform:    "srv",
		EmitterRequestType: "GET",
		EmitterProtocol:    "HTTP",
	}
	ctxR := Context{
		CollectorURI:       "com.acme.override",
		TrackerAppID:       "",
		TrackerNamespace:   "",
		TrackerPlatform:    "",
		EmitterRequestType: "",
		EmitterProtocol:    "",
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
	ctx := Context{
		CollectorURI:       "",
		TrackerAppID:       "",
		TrackerNamespace:   "",
		TrackerPlatform:    "srv",
		EmitterRequestType: "GET",
		EmitterProtocol:    "HTTP",
	}
	ctxR := Context{
		CollectorURI:       "",
		TrackerAppID:       "",
		TrackerNamespace:   "",
		TrackerPlatform:    "",
		EmitterRequestType: "",
		EmitterProtocol:    "",
	}

	trackerChan := make(chan int, 1)
	tracker, err := InitTracker(ctx, ctxR, trackerChan)
	assert.Nil(tracker)
	assert.NotNil(err)
	assert.Equal("URI of the Snowplow Collector is empty - this can be set either at the provider or resource level with the 'collector_uri' input", err.Error())
}

func TestProvider(t *testing.T) {
	assert := assert.New(t)

	p := Provider()
	assert.Nil(p.InternalValidate())
	assert.Contains(p.ResourcesMap, "snowplow_track_self_describing_event")
	assert.Empty(p.DataSourcesMap)
}

func TestProviderConfigure_Defaults(t *testing.T) {
	assert := assert.New(t)

	d := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{})

	ctx, err := providerConfigure(d)
	assert.Nil(err)
	assert.Equal(&Context{
		CollectorURI:       "",
		TrackerAppID:       "",
		TrackerNamespace:   "",
		TrackerPlatform:    "srv",
		EmitterRequestType: "POST",
		EmitterProtocol:    "HTTPS",
	}, ctx)
}

func TestProviderConfigure_WithValues(t *testing.T) {
	assert := assert.New(t)

	d := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{
		"collector_uri":        "com.acme",
		"tracker_app_id":       "app-id",
		"tracker_namespace":    "namespace",
		"tracker_platform":     "web",
		"emitter_request_type": "GET",
		"emitter_protocol":     "HTTP",
	})

	ctx, err := providerConfigure(d)
	assert.Nil(err)
	assert.Equal(&Context{
		CollectorURI:       "com.acme",
		TrackerAppID:       "app-id",
		TrackerNamespace:   "namespace",
		TrackerPlatform:    "web",
		EmitterRequestType: "GET",
		EmitterProtocol:    "HTTP",
	}, ctx)
}
