//
// Copyright (c) 2019-2020 Snowplow Analytics Ltd. All rights reserved.
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
	"fmt"
	gt "github.com/snowplow/snowplow-golang-tracker/v2/tracker"
	"github.com/twinj/uuid"
	"strconv"
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

// contextsFromList converts a list of interfaces to context SDJs.
func contextsFromList(vs []interface{}) ([]gt.SelfDescribingJson, error) {
	result := make([]gt.SelfDescribingJson, 0, len(vs))
	for _, context := range vs {
		attr, ok := context.(map[string]interface{})
		if !ok {
			continue
		}

		t, err := selfDescribingJSONFromMap(attr)
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
func selfDescribingJSONFromMap(attr map[string]interface{}) (*gt.SelfDescribingJson, error) {
	if _, ok := attr["iglu_uri"]; !ok {
		return nil, fmt.Errorf("Invalid attributes: 'iglu_uri' key missing")
	}

	if _, ok := attr["payload"]; !ok {
		return nil, fmt.Errorf("Invalid attributes: 'payload' key missing")
	}

	payload := attr["payload"].(map[string]interface{})
	payloadCopy := make(map[string]interface{})

	for k, v := range payload {
		// Attempt type assertions
		maybeFloat, err := strconv.ParseFloat(v.(string), 64)
		if err == nil {
			payloadCopy[k] = maybeFloat
			continue
		}

		maybeInt, err := strconv.Atoi(v.(string))
		if err == nil {
			payloadCopy[k] = maybeInt
			continue
		}

		maybeBool, err := strconv.ParseBool(v.(string))
		if err == nil {
			payloadCopy[k] = maybeBool
			continue
		}

		// Pass the default value in
		payloadCopy[k] = v
	}

	sdj := gt.InitSelfDescribingJson(
		attr["iglu_uri"].(string),
		payloadCopy,
	)

	return sdj, nil
}
