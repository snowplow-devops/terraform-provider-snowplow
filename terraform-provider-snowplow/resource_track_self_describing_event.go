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
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gt "github.com/snowplow/snowplow-golang-tracker/v2/tracker"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TrackSelfDescribingEventResource{}
var _ resource.ResourceWithConfigure = &TrackSelfDescribingEventResource{}

type TrackSelfDescribingEventResource struct {
	providerContext *SnowplowProviderModel
}

type TrackSelfDescribingEventResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	CreateEvent        types.Map    `tfsdk:"create_event"`
	UpdateEvent        types.Map    `tfsdk:"update_event"`
	DeleteEvent        types.Map    `tfsdk:"delete_event"`
	Contexts           types.List   `tfsdk:"contexts"`
	CollectorURI       types.String `tfsdk:"collector_uri"`
	TrackerAppID       types.String `tfsdk:"tracker_app_id"`
	TrackerNamespace   types.String `tfsdk:"tracker_namespace"`
	TrackerPlatform    types.String `tfsdk:"tracker_platform"`
	EmitterRequestType types.String `tfsdk:"emitter_request_type"`
	EmitterProtocol    types.String `tfsdk:"emitter_protocol"`
}

func NewTrackSelfDescribingEventResource() resource.Resource {
	return &TrackSelfDescribingEventResource{}
}

func (r *TrackSelfDescribingEventResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_track_self_describing_event"
}

func (r *TrackSelfDescribingEventResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Emits an event to the configured collector upon creation, update, or deletion of the resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"create_event": schema.MapAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"update_event": schema.MapAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"delete_event": schema.MapAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"contexts": schema.ListAttribute{
				ElementType: types.MapType{
					ElemType: types.StringType,
				},
				Required: true,
			},
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
				Computed:    true,
				Description: "Optional platform",
				Default:     stringdefault.StaticString("srv"),
			},
			"emitter_request_type": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Computed:    true,
				Description: "Whether to use GET or POST requests to emit events",
				Default:     stringdefault.StaticString("POST"),
			},
			"emitter_protocol": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Computed:    true,
				Description: "Whether to use HTTP or HTTPS to send events",
				Default:     stringdefault.StaticString("HTTPS"),
			},
		},
	}
}

func (r *TrackSelfDescribingEventResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerCtx, ok := req.ProviderData.(SnowplowProviderModel)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected SnowplowProviderModel, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.providerContext = &providerCtx
}

func (r *TrackSelfDescribingEventResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TrackSelfDescribingEventResourceModel

	if resp.Diagnostics.Append((req.Plan.Get(ctx, &plan))...); resp.Diagnostics.HasError() {
		return
	}

	if err := trackSelfDescribingEvent(r.providerContext, plan, plan.CreateEvent); err != nil {
		resp.Diagnostics.AddError("Error tracking event", err.Error())
		return
	}

	if resp.Diagnostics.Append((resp.State.Set(ctx, plan))...); resp.Diagnostics.HasError() {
		return
	}

	resp.State.SetAttribute(ctx, path.Root("id"), getUUID())
}

func (r *TrackSelfDescribingEventResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TrackSelfDescribingEventResourceModel

	resp.Diagnostics.Append((req.State.Get(ctx, &state))...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := trackSelfDescribingEvent(r.providerContext, state, state.DeleteEvent)

	if err != nil {
		resp.Diagnostics.AddError("Error tracking event", err.Error())
		return
	}
}

func (r *TrackSelfDescribingEventResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *TrackSelfDescribingEventResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TrackSelfDescribingEventResourceModel

	if resp.Diagnostics.Append((req.Plan.Get(ctx, &plan))...); resp.Diagnostics.HasError() {
		return
	}

	if err := trackSelfDescribingEvent(r.providerContext, plan, plan.UpdateEvent); err != nil {
		resp.Diagnostics.AddError("Error tracking event", err.Error())
		return
	}

	if resp.Diagnostics.Append((resp.State.Set(ctx, plan))...); resp.Diagnostics.HasError() {
		return
	}

	resp.State.SetAttribute(ctx, path.Root("id"), getUUID())
}

func trackSelfDescribingEvent(providerCtx *SnowplowProviderModel, resourceCtx TrackSelfDescribingEventResourceModel, lifecycleEventMap types.Map) error {
	trackerChan := make(chan int, 1)
	tracker, err := InitTracker(*providerCtx, resourceCtx, trackerChan)
	if err != nil {
		return err
	}

	var contextList = make([]types.Map, 0, len(resourceCtx.Contexts.Elements()))

	for _, elem := range resourceCtx.Contexts.Elements() {
		if entity, ok := elem.(types.Map); ok {
			contextList = append(contextList, entity)
		} else {
			return errors.New("All contexts values must be maps")
		}
	}

	contexts, err := contextsFromList(contextList)
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
