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
	"github.com/stretchr/testify/assert"
	"testing"
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
