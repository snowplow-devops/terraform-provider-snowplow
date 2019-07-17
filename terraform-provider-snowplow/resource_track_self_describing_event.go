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

func resourceTrackSelfDescribingEvent() *schema.Resource {
	return &schema.Resource{
		Create: resourceTrackSelfDescribingEventCreate,
		Read:   resourceTrackSelfDescribingEventRead,
		Update: resourceTrackSelfDescribingEventUpdate,
		Delete: resourceTrackSelfDescribingEventDelete,

		Schema: map[string]*schema.Schema{
			"iglu_uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"payload": {
				Type:     schema.TypeString,
				Required: true,
			},
			"contexts": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"create_context": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"update_context": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"delete_context": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func trackSelfDescribingEvent(d *schema.ResourceData, m interface{}, lifecycleContextMap map[string]interface{}) error {
	ctx := m.(*Context)

	trackerChan := make(chan int, 1)
	tracker := InitTracker(*ctx, trackerChan)

	contexts, err := contextsFromList(d.Get("contexts").([]interface{}))
	if err != nil {
		return err
	}

	if lifecycleContextMap != nil {
		lifecycleContext, err := contextFromMap(lifecycleContextMap)
		if err != nil {
			return err
		}

		if lifecycleContext != nil {
			contexts = append(contexts, *lifecycleContext)
		}
	}

	payloadData, err := stringToMap(d.Get("payload").(string))
	if err != nil {
		return err
	}

	igluURI := d.Get("iglu_uri").(string)

	sdj := gt.InitSelfDescribingJson(
		igluURI,
		payloadData,
	)

	tracker.TrackSelfDescribingEvent(gt.SelfDescribingEvent{
		Event:    sdj,
		Contexts: contexts,
	})

	statusCode := <-trackerChan

	return parseStatusCode(statusCode)
}

func resourceTrackSelfDescribingEventCreate(d *schema.ResourceData, m interface{}) error {
	err := trackSelfDescribingEvent(d, m, d.Get("create_context").(map[string]interface{}))
	if err != nil {
		return err
	}

	d.SetId(d.Get("iglu_uri").(string))

	return resourceTrackSelfDescribingEventRead(d, m)
}

func resourceTrackSelfDescribingEventRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceTrackSelfDescribingEventUpdate(d *schema.ResourceData, m interface{}) error {
	err := trackSelfDescribingEvent(d, m, d.Get("update_context").(map[string]interface{}))
	if err != nil {
		return err
	}

	return resourceTrackSelfDescribingEventRead(d, m)
}

func resourceTrackSelfDescribingEventDelete(d *schema.ResourceData, m interface{}) error {
	err := trackSelfDescribingEvent(d, m, d.Get("delete_context").(map[string]interface{}))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
