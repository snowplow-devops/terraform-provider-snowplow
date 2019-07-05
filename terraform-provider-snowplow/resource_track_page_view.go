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
	"fmt"
)

func resourceTrackPageView() *schema.Resource {
	return &schema.Resource{
		Create: resourceTrackPageViewCreate,
		Read:   resourceTrackPageViewRead,
		Delete: resourceTrackPageViewDelete,

		Schema: map[string]*schema.Schema{
			"pv_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"pv_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceTrackPageViewCreate(d *schema.ResourceData, m interface{}) error {
	ctx := m.(*Context)

	trackerChan := make(chan int, 1)
	tracker := InitTracker(*ctx, trackerChan)

	tracker.TrackScreenView(gt.ScreenViewEvent{
		Name: gt.NewString(d.Get("pv_name").(string)),
		Id:   gt.NewString(d.Get("pv_id").(string)),
	})

	statusCode := <-trackerChan

	if !ParseStatusCode(statusCode) {
		return fmt.Errorf("Got %d status code when sending event - need 2xx or 3xx", statusCode)
	}

	return resourceTrackPageViewRead(d, m)
}

func resourceTrackPageViewRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceTrackPageViewDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
