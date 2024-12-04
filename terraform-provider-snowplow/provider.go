//
// Copyright (c) 2019-2023 Snowplow Analytics Ltd. All rights reserved.
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
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	gt "github.com/snowplow/snowplow-golang-tracker/v2/tracker"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &SnowplowProvider{}

// SnowplowProviderModel the struct made from the provider input options
type SnowplowProviderModel struct {
	CollectorURI       types.String `tfsdk:"collector_uri"`
	TrackerAppID       types.String `tfsdk:"tracker_app_id"`
	TrackerNamespace   types.String `tfsdk:"tracker_namespace"`
	TrackerPlatform    types.String `tfsdk:"tracker_platform"`
	EmitterRequestType types.String `tfsdk:"emitter_request_type"`
	EmitterProtocol    types.String `tfsdk:"emitter_protocol"`
}

// SnowplowProvider defines the provider implementation.
type SnowplowProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func NewProvider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SnowplowProvider{version: version}
	}
}

func (p *SnowplowProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "snowplow"
	resp.Version = p.version
}

func (p *SnowplowProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for emitting Snowplow events",
		Attributes: map[string]schema.Attribute{
			"collector_uri": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Description: "URI of your Snowplow Collector",
			},
			"tracker_app_id": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Description: "Optional application ID",
			},
			"tracker_namespace": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Description: "Optional namespace",
			},
			"tracker_platform": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Description: "Optional platform",
			},
			"emitter_request_type": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Description: "Whether to use GET or POST requests to emit events",
			},
			"emitter_protocol": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Description: "Whether to use HTTP or HTTPS to send events",
			},
		},
	}
}

func (p *SnowplowProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data SnowplowProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = data
	resp.ResourceData = data
}

func (p *SnowplowProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTrackSelfDescribingEventResource,
	}
}

func (p *SnowplowProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

// InitTracker takes a context and a channel of size 1 and returns
// a new Snowplow Tracker ready to create a resource
func InitTracker(ctx SnowplowProviderModel, ctxResource TrackSelfDescribingEventResourceModel, trackerChan chan int) (*gt.Tracker, error) {
	var collectorUri, emitterRequestType, emitterProtocol, trackerNamespace, trackerAppId, trackerPlatform string

	if ctxResource.CollectorURI.IsNull() || ctxResource.CollectorURI.ValueString() == "" {
		collectorUri = ctx.CollectorURI.ValueString()
	} else {
		collectorUri = ctxResource.CollectorURI.ValueString()
	}

	if ctxResource.EmitterRequestType.IsNull() || ctxResource.EmitterRequestType.ValueString() == "" {
		emitterRequestType = ctx.EmitterRequestType.ValueString()
	} else {
		emitterRequestType = ctxResource.EmitterRequestType.ValueString()
	}

	if ctxResource.EmitterProtocol.IsNull() || ctxResource.EmitterProtocol.ValueString() == "" {
		emitterProtocol = ctx.EmitterProtocol.ValueString()
	} else {
		emitterProtocol = ctxResource.EmitterProtocol.ValueString()
	}

	if ctxResource.TrackerNamespace.IsNull() || ctxResource.TrackerNamespace.ValueString() == "" {
		trackerNamespace = ctx.TrackerNamespace.ValueString()
	} else {
		trackerNamespace = ctxResource.TrackerNamespace.ValueString()
	}

	if ctxResource.TrackerAppID.IsNull() || ctxResource.TrackerAppID.ValueString() == "" {
		trackerAppId = ctx.TrackerAppID.ValueString()
	} else {
		trackerAppId = ctxResource.TrackerAppID.ValueString()
	}

	if ctxResource.TrackerPlatform.IsNull() || ctxResource.TrackerPlatform.ValueString() == "" {
		trackerPlatform = ctx.TrackerPlatform.ValueString()
	} else {
		trackerPlatform = ctxResource.TrackerPlatform.ValueString()
	}

	if collectorUri == "" {
		return nil, errors.New("URI of the Snowplow Collector is empty - this can be set either at the provider or resource level with the 'collector_uri' input")
	}

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
		gt.RequireCollectorUri(collectorUri),
		gt.OptionRequestType(emitterRequestType),
		gt.OptionProtocol(emitterProtocol),
		gt.OptionCallback(callback),
		gt.OptionStorage(gt.InitStorageMemory()),
	)

	subject := gt.InitSubject()

	tracker := gt.InitTracker(
		gt.RequireEmitter(emitter),
		gt.OptionSubject(subject),
		gt.OptionNamespace(trackerNamespace),
		gt.OptionAppId(trackerAppId),
		gt.OptionPlatform(trackerPlatform),
		gt.OptionBase64Encode(true),
	)

	return tracker, nil
}
