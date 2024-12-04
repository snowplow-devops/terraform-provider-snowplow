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
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	gt "github.com/snowplow/snowplow-golang-tracker/v2/tracker"
	"github.com/twinj/uuid"
)

// getUUID generates a Version 4 UUID string.
func getUUID() string {
	return uuid.NewV4().String()
}

// parseStatusCode checks whether we got a valid status code from
// the collector.
func parseStatusCode(statusCode int) error {
	var err error
	result := statusCode / 100

	switch result {
	case 2, 3:
		err = nil
	default:
		err = fmt.Errorf("Got %d status code when sending event - need 2xx or 3xx", statusCode)
	}

	return err
}

// stringToMap attempts to convert a string (assumed JSON) to a map.
func stringToMap(str string) (map[string]interface{}, error) {
	var jsonDataMap map[string]interface{}
	d := json.NewDecoder(strings.NewReader(str))
	d.UseNumber()
	err := d.Decode(&jsonDataMap)
	if err != nil {
		return nil, err
	}
	return jsonDataMap, nil
}

// contextsFromList converts a list of interfaces to context SDJs.
func contextsFromList(vs []types.Map) ([]gt.SelfDescribingJson, error) {
	result := make([]gt.SelfDescribingJson, 0, len(vs))
	for _, context := range vs {
		t, err := selfDescribingJSONFromMap(context)
		if err != nil {
			return nil, err
		}

		if t != nil {
			result = append(result, *t)
		}
	}
	return result, nil
}

// selfDescribingJsonFromMap converts a map into a context SDJ.
func selfDescribingJSONFromMap(obj types.Map) (*gt.SelfDescribingJson, error) {
	attr := obj.Elements()

	var igluUri, payload string

	if val, ok := attr["iglu_uri"]; !ok {
		return nil, fmt.Errorf("Invalid context attributes: 'iglu_uri' key missing")
	} else if igluUriVal, ok := val.(types.String); !ok {
		return nil, fmt.Errorf("Invalid context attributes: 'iglu_uri' not string")
	} else {
		igluUri = igluUriVal.ValueString()
	}

	if val, ok := attr["payload"]; !ok {
		return nil, fmt.Errorf("Invalid context attributes: 'payload' key missing")
	} else if payloadVal, ok := val.(types.String); !ok {
		return nil, fmt.Errorf("Invalid context attributes: 'payload' not string")
	} else {
		payload = payloadVal.ValueString()
	}

	contextData, err := stringToMap(payload)
	if err != nil {
		return nil, err
	}

	sdj := gt.InitSelfDescribingJson(
		igluUri,
		contextData,
	)

	return sdj, nil
}
