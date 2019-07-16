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
		Delete: resourceTrackSelfDescribingEventDelete,

		Schema: map[string]*schema.Schema{
			"iglu_uri": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"payload": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"contexts": {
				Type:     schema.TypeList,
				ForceNew: true,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
		},
	}
}

func resourceTrackSelfDescribingEventCreate(d *schema.ResourceData, m interface{}) error {
	ctx := m.(*Context)

	trackerChan := make(chan int, 1)
	tracker := InitTracker(*ctx, trackerChan)

	contexts, err := contextsFromList(d.Get("contexts").([]interface{}))
	if err != nil {
		return err
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

	err = parseStatusCode(statusCode)
	if err != nil {
		return err
	}

	d.SetId(igluURI)

	return resourceTrackSelfDescribingEventRead(d, m)
}

func resourceTrackSelfDescribingEventRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceTrackSelfDescribingEventDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}
