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
	gt "gopkg.in/snowplow/snowplow-golang-tracker.v2/tracker"
)

type Context struct {
	tracker *gt.Tracker
}

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
			"snowplow_track_page_view": resourceTrackPageView(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureFunc:  providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	emitter := gt.InitEmitter(
		gt.RequireCollectorUri(d.Get("collector_uri").(string)),
		gt.OptionRequestType(d.Get("emitter_request_type").(string)),
		gt.OptionProtocol(d.Get("emitter_protocol").(string)),
		gt.OptionStorage(gt.InitStorageMemory()),
	)

	subject := gt.InitSubject()

	tracker := gt.InitTracker(
		gt.RequireEmitter(emitter),
		gt.OptionSubject(subject),
		gt.OptionNamespace(d.Get("tracker_namespace").(string)),
		gt.OptionAppId(d.Get("tracker_app_id").(string)),
		gt.OptionPlatform(d.Get("tracker_platform").(string)),
		gt.OptionBase64Encode(true),
	)

	ctx := Context{
		tracker: tracker,
	}

	return &ctx, nil
}
