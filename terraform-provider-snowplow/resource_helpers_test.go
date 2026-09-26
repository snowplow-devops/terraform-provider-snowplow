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
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseStatusCode_2xx3xx(t *testing.T) {
	assert := assert.New(t)

	err := parseStatusCode(200)
	assert.Nil(err)

	err = parseStatusCode(300)
	assert.Nil(err)
}

func TestParseStatusCode_4xx5xx(t *testing.T) {
	assert := assert.New(t)

	err := parseStatusCode(404)
	assert.NotNil(err)

	err = parseStatusCode(504)
	assert.NotNil(err)
}

func TestStringToMap(t *testing.T) {
	assert := assert.New(t)

	m, err := stringToMap("{\"hello\":\"world\"}")
	assert.Nil(err)
	assert.NotNil(m)
	assert.Equal("world", m["hello"])
	assert.Equal(1, len(m))

	m, err = stringToMap("{\"hello\"}")
	assert.NotNil(err)
	assert.Nil(m)

	m, err = stringToMap("{\"timestamp\":1534429336}")
	assert.Nil(err)
	assert.NotNil(m)
	assert.Equal(json.Number("1534429336"), m["timestamp"])
	assert.Equal(1, len(m))
}

func TestContextsFromList_Valid(t *testing.T) {
	assert := assert.New(t)

	contextList := make([]interface{}, 0, 2)
	contextList = append(contextList, map[string]interface{}{"iglu_uri": "iglu:com.acme/context_1/jsonschema/1-0-0", "payload": "{\"foo\":\"bar\"}"})
	contextList = append(contextList, map[string]interface{}{"iglu_uri": "iglu:com.acme/context_2/jsonschema/1-0-0", "payload": "{\"foo2\":\"bar2\"}"})

	contextSdeList, err := contextsFromList(contextList)
	assert.Nil(err)
	assert.NotNil(contextSdeList)
	assert.Equal(2, len(contextSdeList))
}

func TestContextsFromList_NoIgluUri(t *testing.T) {
	assert := assert.New(t)

	contextList := make([]interface{}, 0, 1)
	contextList = append(contextList, map[string]interface{}{"iglu_uriss": "iglu:com.acme/context_1/jsonschema/1-0-0", "payload": "{\"foo\":\"bar\"}"})

	contextSdeList, err := contextsFromList(contextList)
	assert.NotNil(err)
	assert.Nil(contextSdeList)
}

func TestContextsFromList_NoPayload(t *testing.T) {
	assert := assert.New(t)

	contextList := make([]interface{}, 0, 1)
	contextList = append(contextList, map[string]interface{}{"iglu_uri": "iglu:com.acme/context_1/jsonschema/1-0-0", "payloadsss": "{\"foo\":\"bar\"}"})

	contextSdeList, err := contextsFromList(contextList)
	assert.NotNil(err)
	assert.Nil(contextSdeList)
}

func TestContextsFromList_SkipsNonMaps(t *testing.T) {
	assert := assert.New(t)

	contextList := []interface{}{
		nil,
		"not-a-map",
		map[string]interface{}{"iglu_uri": "iglu:com.acme/context_1/jsonschema/1-0-0", "payload": "{\"foo\":\"bar\"}"},
	}

	contextSdeList, err := contextsFromList(contextList)
	assert.Nil(err)
	assert.Equal(1, len(contextSdeList))
}

func TestContextsFromList_Empty(t *testing.T) {
	assert := assert.New(t)

	contextSdeList, err := contextsFromList([]interface{}{})
	assert.Nil(err)
	assert.NotNil(contextSdeList)
	assert.Equal(0, len(contextSdeList))
}

func TestSelfDescribingJSONFromMap_Valid(t *testing.T) {
	assert := assert.New(t)

	sdj, err := selfDescribingJSONFromMap(map[string]interface{}{"iglu_uri": "iglu:com.acme/event/jsonschema/1-0-0", "payload": "{\"foo\":\"bar\"}"})
	assert.Nil(err)
	assert.NotNil(sdj)
	assert.Equal("{\"data\":{\"foo\":\"bar\"},\"schema\":\"iglu:com.acme/event/jsonschema/1-0-0\"}", sdj.String())
}

func TestSelfDescribingJSONFromMap_InvalidPayload(t *testing.T) {
	assert := assert.New(t)

	sdj, err := selfDescribingJSONFromMap(map[string]interface{}{"iglu_uri": "iglu:com.acme/event/jsonschema/1-0-0", "payload": "{\"foo\"}"})
	assert.NotNil(err)
	assert.Nil(sdj)
}

func TestGetUUID(t *testing.T) {
	assert := assert.New(t)

	id := getUUID()
	assert.Regexp("^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$", id)
	assert.NotEqual(id, getUUID())
}
