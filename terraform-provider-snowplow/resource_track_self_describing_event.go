//
// Copyright (c) 2019-2020 Snowplow Analytics Ltd. All rights reserved.
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

func resourceTrackSelfDescribingEvent() *schema.Resource {
	return &schema.Resource{
		Create: resourceTrackSelfDescribingEventCreate,
		Read:   resourceTrackSelfDescribingEventRead,
		Update: resourceTrackSelfDescribingEventUpdate,
		Delete: resourceTrackSelfDescribingEventDelete,

		Schema: map[string]*schema.Schema{
			"create_event": {
				Type:     schema.TypeMap,
				Required: true,
			},
			"update_event": {
				Type:     schema.TypeMap,
				Required: true,
			},
			"delete_event": {
				Type:     schema.TypeMap,
				Required: true,
			},
			"contexts": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"collector_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Required:    false,
				Description: "URI of your Snowplow Collector",
				Default:     "",
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
				Default:     "",
			},
			"emitter_request_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Required:    false,
				Description: "Whether to use GET or POST requests to emit events",
				Default:     "",
			},
			"emitter_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Required:    false,
				Description: "Whether to use HTTP or HTTPS to send events",
				Default:     "",
			},
		},
	}
}

func trackSelfDescribingEvent(d *schema.ResourceData, m interface{}, lifecycleEventMap map[string]interface{}) error {
	ctx := m.(*Context)

	ctxResource, err := providerConfigure(d)
	if err != nil {
		return err
	}
	ctxResourceCast := ctxResource.(*Context)

	trackerChan := make(chan int, 1)
	tracker, err := InitTracker(*ctx, *ctxResourceCast, trackerChan)
	if err != nil {
		return err
	}

	contexts, err := contextsFromList(d.Get("contexts").([]interface{}))
	if err != nil {
		return err
	}

	lifecycleEvent, err := selfDescribingJSONFromMap(lifecycleEventMap)
	if err != nil {
		return err
	}

	tracker.TrackSelfDescribingEvent(gt.SelfDescribingEvent{
		Event:    lifecycleEvent,
		Contexts: contexts,
	})

	statusCode := <-trackerChan

	return parseStatusCode(statusCode)
}

func resourceTrackSelfDescribingEventCreate(d *schema.ResourceData, m interface{}) error {
	err := trackSelfDescribingEvent(d, m, d.Get("create_event").(map[string]interface{}))
	if err != nil {
		return err
	}

	d.SetId(getUUID())

	return resourceTrackSelfDescribingEventRead(d, m)
}

func resourceTrackSelfDescribingEventRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceTrackSelfDescribingEventUpdate(d *schema.ResourceData, m interface{}) error {
	err := trackSelfDescribingEvent(d, m, d.Get("update_event").(map[string]interface{}))
	if err != nil {
		return err
	}

	return resourceTrackSelfDescribingEventRead(d, m)
}

func resourceTrackSelfDescribingEventDelete(d *schema.ResourceData, m interface{}) error {
	err := trackSelfDescribingEvent(d, m, d.Get("delete_event").(map[string]interface{}))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
