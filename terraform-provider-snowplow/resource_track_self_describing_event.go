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
	eventSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"iglu_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payload": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}

	return &schema.Resource{
		Create: resourceTrackSelfDescribingEventCreate,
		Read:   resourceTrackSelfDescribingEventRead,
		Update: resourceTrackSelfDescribingEventUpdate,
		Delete: resourceTrackSelfDescribingEventDelete,

		Schema: map[string]*schema.Schema{
			"create_event": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     eventSchema,
			},
			"update_event": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     eventSchema,
			},
			"delete_event": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     eventSchema,
			},
			"context": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     eventSchema,
			},
		},
	}
}

func trackSelfDescribingEvent(d *schema.ResourceData, m interface{}, lifecycleEventMap interface{}) error {
	ctx := m.(*Context)

	trackerChan := make(chan int, 1)
	tracker := InitTracker(*ctx, trackerChan)

	contexts, err := contextsFromList(d.Get("context").([]interface{}))
	if err != nil {
		return err
	}

	lifecycleEvent, err := selfDescribingJSONFromMap(lifecycleEventMap.(map[string]interface{}))
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
	err := trackSelfDescribingEvent(d, m, d.Get("create_event").([]interface{})[0])
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
	err := trackSelfDescribingEvent(d, m, d.Get("update_event").([]interface{})[0])
	if err != nil {
		return err
	}

	return resourceTrackSelfDescribingEventRead(d, m)
}

func resourceTrackSelfDescribingEventDelete(d *schema.ResourceData, m interface{}) error {
	err := trackSelfDescribingEvent(d, m, d.Get("delete_event").([]interface{})[0])
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
