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
	"github.com/hashicorp/terraform/helper/schema"
	gt "github.com/snowplow/snowplow-golang-tracker/v2/tracker"
)

// Context the struct made from the provider input options
type Context struct {
	CollectorURI       string
	TrackerAppID       string
	TrackerNamespace   string
	TrackerPlatform    string
	EmitterRequestType string
	EmitterProtocol    string
}

// Provider creates a new provider struct ready for use by Terraform
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"collector_uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URI of your Snowplow Collector",
			},
			"tracker_app_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Required:    false,
				Description: "Optional application ID",
				Default:     "",
			},
			"tracker_namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Required:    false,
				Description: "Optional namespace",
				Default:     "",
			},
			"tracker_platform": {
				Type:        schema.TypeString,
				Optional:    true,
				Required:    false,
				Description: "Optional platform",
				Default:     "srv",
			},
			"emitter_request_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Required:    false,
				Description: "Whether to use GET or POST requests to emit events",
				Default:     "POST",
			},
			"emitter_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Required:    false,
				Description: "Whether to use HTTP or HTTPS to send events",
				Default:     "HTTPS",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"snowplow_track_self_describing_event": resourceTrackSelfDescribingEvent(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureFunc:  providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	ctx := Context{
		CollectorURI:       d.Get("collector_uri").(string),
		TrackerAppID:       d.Get("tracker_app_id").(string),
		TrackerNamespace:   d.Get("tracker_namespace").(string),
		TrackerPlatform:    d.Get("tracker_platform").(string),
		EmitterRequestType: d.Get("emitter_request_type").(string),
		EmitterProtocol:    d.Get("emitter_protocol").(string),
	}

	return &ctx, nil
}

// InitTracker takes a context and a channel of size 1 and returns
// a new Snowplow Tracker ready to create a resource
func InitTracker(ctx Context, trackerChan chan int) *gt.Tracker {
	callback := func(s []gt.CallbackResult, f []gt.CallbackResult) {
		status := 0

		if len(s) == 1 {
			status = s[0].Status
		} else if len(f) == 1 {
			status = f[0].Status
		}

		trackerChan <- status
	}

	emitter := gt.InitEmitter(
		gt.RequireCollectorUri(ctx.CollectorURI),
		gt.OptionRequestType(ctx.EmitterRequestType),
		gt.OptionProtocol(ctx.EmitterProtocol),
		gt.OptionCallback(callback),
		gt.OptionStorage(gt.InitStorageMemory()),
	)

	subject := gt.InitSubject()

	tracker := gt.InitTracker(
		gt.RequireEmitter(emitter),
		gt.OptionSubject(subject),
		gt.OptionNamespace(ctx.TrackerNamespace),
		gt.OptionAppId(ctx.TrackerAppID),
		gt.OptionPlatform(ctx.TrackerPlatform),
		gt.OptionBase64Encode(true),
	)

	return tracker
}
