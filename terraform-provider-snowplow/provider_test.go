//
// Copyright (c) 2019 Snowplow Analytics Ltd. All rights reserved.
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
		CollectorUri:       "com.acme",
		TrackerAppId:       "",
		TrackerNamespace:   "",
		TrackerPlatform:    "srv",
		EmitterRequestType: "GET",
		EmitterProtocol:    "HTTP",
	}

	trackerChan := make(chan int, 1)
	tracker := InitTracker(ctx, trackerChan)
	assert.NotNil(tracker)
}
