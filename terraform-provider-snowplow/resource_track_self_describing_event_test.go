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
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// mockCollector is a fake Snowplow Collector which records the requests
// it receives and responds with a fixed status code.
type mockCollector struct {
	server   *httptest.Server
	mu       sync.Mutex
	requests []mockRequest
}

type mockRequest struct {
	method string
	query  url.Values
	body   string
}

func newMockCollector(t *testing.T, statusCode int) *mockCollector {
	c := &mockCollector{}
	c.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		c.mu.Lock()
		c.requests = append(c.requests, mockRequest{method: r.Method, query: r.URL.Query(), body: string(body)})
		c.mu.Unlock()
		w.WriteHeader(statusCode)
	}))
	t.Cleanup(c.server.Close)
	return c
}

// host returns the collector address without the scheme, as expected
// by the 'collector_uri' input.
func (c *mockCollector) host() string {
	return strings.TrimPrefix(c.server.URL, "http://")
}

func (c *mockCollector) requestCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.requests)
}

func (c *mockCollector) lastRequest() mockRequest {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.requests[len(c.requests)-1]
}

func providerContext(collectorURI string) *Context {
	return &Context{
		CollectorURI:       collectorURI,
		TrackerAppID:       "app-id",
		TrackerNamespace:   "namespace",
		TrackerPlatform:    "srv",
		EmitterRequestType: "POST",
		EmitterProtocol:    "HTTP",
	}
}

func lifecycleEvent(name string) map[string]interface{} {
	return map[string]interface{}{
		"iglu_uri": "iglu:com.acme/lifecycle/jsonschema/1-0-0",
		"payload":  "{\"action\":\"" + name + "\"}",
	}
}

func resourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	base := map[string]interface{}{
		"create_event": lifecycleEvent("create"),
		"update_event": lifecycleEvent("update"),
		"delete_event": lifecycleEvent("delete"),
		"contexts": []interface{}{
			map[string]interface{}{
				"iglu_uri": "iglu:com.acme/context/jsonschema/1-0-0",
				"payload":  "{\"foo\":\"bar\"}",
			},
		},
	}
	for k, v := range raw {
		base[k] = v
	}
	return schema.TestResourceDataRaw(t, resourceTrackSelfDescribingEvent().Schema, base)
}

func TestResourceTrackSelfDescribingEvent_Schema(t *testing.T) {
	assert := assert.New(t)

	r := resourceTrackSelfDescribingEvent()
	assert.Nil(r.InternalValidate(nil, true))

	for _, k := range []string{"create_event", "update_event", "delete_event", "contexts"} {
		assert.True(r.Schema[k].Required, k)
	}
	for _, k := range []string{"collector_uri", "tracker_app_id", "tracker_namespace", "tracker_platform", "emitter_request_type", "emitter_protocol"} {
		assert.True(r.Schema[k].Optional, k)
		assert.Equal("", r.Schema[k].Default, k)
	}
}

func TestResourceTrackSelfDescribingEvent_Create(t *testing.T) {
	assert := assert.New(t)

	collector := newMockCollector(t, http.StatusOK)
	d := resourceData(t, nil)

	err := resourceTrackSelfDescribingEventCreate(d, providerContext(collector.host()))
	assert.Nil(err)
	assert.Len(d.Id(), 36)
	assert.Equal(1, collector.requestCount())
	req := collector.lastRequest()
	assert.Equal(http.MethodPost, req.method)
	assert.Contains(req.body, "\"aid\":\"app-id\"")
	assert.Contains(req.body, "\"tna\":\"namespace\"")
	assert.Contains(req.body, "\"p\":\"srv\"")
}

func TestResourceTrackSelfDescribingEvent_Update(t *testing.T) {
	assert := assert.New(t)

	collector := newMockCollector(t, http.StatusOK)
	d := resourceData(t, nil)
	d.SetId("existing-id")

	err := resourceTrackSelfDescribingEventUpdate(d, providerContext(collector.host()))
	assert.Nil(err)
	assert.Equal("existing-id", d.Id())
	assert.Equal(1, collector.requestCount())
}

func TestResourceTrackSelfDescribingEvent_Delete(t *testing.T) {
	assert := assert.New(t)

	collector := newMockCollector(t, http.StatusOK)
	d := resourceData(t, nil)
	d.SetId("existing-id")

	err := resourceTrackSelfDescribingEventDelete(d, providerContext(collector.host()))
	assert.Nil(err)
	assert.Equal("", d.Id())
	assert.Equal(1, collector.requestCount())
}

func TestResourceTrackSelfDescribingEvent_Read(t *testing.T) {
	assert := assert.New(t)

	d := resourceData(t, nil)
	d.SetId("existing-id")

	assert.Nil(resourceTrackSelfDescribingEventRead(d, providerContext("")))
	assert.Equal("existing-id", d.Id())
}

func TestResourceTrackSelfDescribingEvent_ResourceOverridesProvider(t *testing.T) {
	assert := assert.New(t)

	collector := newMockCollector(t, http.StatusOK)
	d := resourceData(t, map[string]interface{}{
		"collector_uri":        collector.host(),
		"tracker_app_id":       "resource-app-id",
		"emitter_request_type": "GET",
	})

	// The provider points at an unreachable collector so the event can only
	// arrive if the resource level 'collector_uri' takes precedence
	err := resourceTrackSelfDescribingEventCreate(d, providerContext("127.0.0.1:1"))
	assert.Nil(err)
	assert.Equal(1, collector.requestCount())

	// Overridden values come from the resource, the rest from the provider
	req := collector.lastRequest()
	assert.Equal(http.MethodGet, req.method)
	assert.Equal("resource-app-id", req.query.Get("aid"))
	assert.Equal("namespace", req.query.Get("tna"))
	assert.Equal("srv", req.query.Get("p"))
}

func TestResourceTrackSelfDescribingEvent_CollectorError(t *testing.T) {
	assert := assert.New(t)

	collector := newMockCollector(t, http.StatusInternalServerError)
	ctx := providerContext(collector.host())

	d := resourceData(t, nil)
	err := resourceTrackSelfDescribingEventCreate(d, ctx)
	assert.EqualError(err, "got 500 status code when sending event - need 2xx or 3xx")
	assert.Equal("", d.Id())

	d = resourceData(t, nil)
	d.SetId("existing-id")
	assert.NotNil(resourceTrackSelfDescribingEventUpdate(d, ctx))

	d = resourceData(t, nil)
	d.SetId("existing-id")
	assert.NotNil(resourceTrackSelfDescribingEventDelete(d, ctx))
	assert.Equal("existing-id", d.Id())
}

func TestResourceTrackSelfDescribingEvent_EmptyCollectorURI(t *testing.T) {
	assert := assert.New(t)

	d := resourceData(t, nil)
	err := resourceTrackSelfDescribingEventCreate(d, providerContext(""))
	assert.ErrorContains(err, "URI of the Snowplow Collector is empty")
	assert.Equal("", d.Id())
}

func TestResourceTrackSelfDescribingEvent_InvalidContext(t *testing.T) {
	assert := assert.New(t)

	collector := newMockCollector(t, http.StatusOK)
	d := resourceData(t, map[string]interface{}{
		"contexts": []interface{}{
			map[string]interface{}{"payload": "{\"foo\":\"bar\"}"},
		},
	})

	err := resourceTrackSelfDescribingEventCreate(d, providerContext(collector.host()))
	assert.EqualError(err, "invalid context attributes: 'iglu_uri' key missing")
	assert.Equal(0, collector.requestCount())
}

func TestResourceTrackSelfDescribingEvent_InvalidLifecycleEvent(t *testing.T) {
	assert := assert.New(t)

	collector := newMockCollector(t, http.StatusOK)
	d := resourceData(t, map[string]interface{}{
		"create_event": map[string]interface{}{"iglu_uri": "iglu:com.acme/lifecycle/jsonschema/1-0-0"},
	})

	err := resourceTrackSelfDescribingEventCreate(d, providerContext(collector.host()))
	assert.EqualError(err, "invalid context attributes: 'payload' key missing")
	assert.Equal(0, collector.requestCount())
}
