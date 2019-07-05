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
				Required:    false,
				Description: "Optional application ID",
				Default:     "",
			},
			"tracker_namespace": {
				Type:        schema.TypeString,
				Required:    false,
				Description: "Optional namespace",
				Default:     "",
			},
			"tracker_platform": {
				Type:        schema.TypeString,
				Required:    false,
				Description: "Optional platform",
				Default:     "srv",
			},
			"base64_encode": {
				Type:        schema.TypeBool,
				Required:    false,
				Description: "Whether to base64 encode custom contexts and self-describing JSONs",
				Default:     true,
			},
			"emitter_request_type": {
				Type:        schema.TypeString,
				Required:    false,
				Description: "Whether to use GET or POST requests to emit events",
				Default:     "POST",
			},
			"emitter_protocol": {
				Type:        schema.TypeString,
				Required:    false,
				Description: "Whether to use HTTP or HTTPS to send events",
				Default:     "HTTPS",
			},
		},
		ResourcesMap:  map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureFunc: providerConfigure,
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
		gt.OptionBase64Encode(d.Get("base64_encode").(bool)),
	)

	ctx := Context{
		tracker: tracker,
	}

	return &ctx, nil
}
