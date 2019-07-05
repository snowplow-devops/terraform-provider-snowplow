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

func resourceTrackPageView() *schema.Resource {
	return &schema.Resource{
		Create: resourceTrackPageViewCreate,
		Read:   resourceTrackPageViewRead,
		Update: resourceTrackPageViewUpdate,
		Delete: resourceTrackPageViewDelete,

		Schema: map[string]*schema.Schema{
			"pv_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pv_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceTrackPageViewCreate(d *schema.ResourceData, m interface{}) error {
	ctx := m.(*Context)
	tracker := ctx.tracker

	tracker.TrackScreenView(gt.ScreenViewEvent{
	  Name: gt.NewString(d.Get("pv_name").(string)),
	  Id: gt.NewString(d.Get("pv_id").(string)),
	})

	return resourceTrackPageViewRead(d, m)
}

func resourceTrackPageViewRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceTrackPageViewUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceTrackPageViewCreate(d, m)
}

func resourceTrackPageViewDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
