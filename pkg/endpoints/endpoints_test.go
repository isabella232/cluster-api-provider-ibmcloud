/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package endpoints

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFlags(t *testing.T) {
	testCases := []struct {
		name           string
		flagToParse    string
		expectedOutput []ServiceEndpoint
		expectedError  error
	}{
		{
			name:           "no configuration",
			flagToParse:    "",
			expectedOutput: nil,
			expectedError:  nil,
		},
		{
			name:        "single region, single vpc service",
			flagToParse: "eu-gb:vpc=https://vpchost:8080",
			expectedOutput: []ServiceEndpoint{
				{
					ID:     "vpc",
					URL:    "https://vpchost:8080",
					Region: "eu-gb",
				},
			},
			expectedError: nil,
		},
		{
			name:        "single region, single powervs service",
			flagToParse: "lon:powervs=https://pvshost:8080",
			expectedOutput: []ServiceEndpoint{
				{
					ID:     "powervs",
					URL:    "https://pvshost:8080",
					Region: "lon",
				},
			},
			expectedError: nil,
		},
		{
			name:        "single region, multiple services",
			flagToParse: "lon:powervs=https://pvshost:8080,rc=https://rchost:8080",
			expectedOutput: []ServiceEndpoint{
				{
					ID:     "powervs",
					URL:    "https://pvshost:8080",
					Region: "lon",
				},
				{
					ID:     "rc",
					URL:    "https://rchost:8080",
					Region: "lon",
				},
			},
			expectedError: nil,
		},
		{
			name:           "single region, duplicate service",
			flagToParse:    "eu-gb:vpc=https://localhost:8080,vpc=https://vpchost:8080",
			expectedOutput: nil,
			expectedError:  errServiceEndpointDuplicateID,
		},
		{
			name:           "single region, non-valid URI",
			flagToParse:    "eu-gb:vpc=fdsfs",
			expectedOutput: nil,
			expectedError:  errServiceEndpointURL,
		},
		{
			name:        "multiples regions",
			flagToParse: "eu-gb:vpc=https://vpchost:8080;lon:powervs=https://pvshost:8080,rc=https://rchost:8080",
			expectedOutput: []ServiceEndpoint{
				{
					ID:     "vpc",
					URL:    "https://vpchost:8080",
					Region: "eu-gb",
				},
				{
					ID:     "powervs",
					URL:    "https://pvshost:8080",
					Region: "lon",
				},
				{
					ID:     "rc",
					URL:    "https://rchost:8080",
					Region: "lon",
				},
			},
			expectedError: nil,
		},
		{
			name:        "multiples regions, multiple services",
			flagToParse: "eu-gb:vpc=https://vpchost:8080;lon:powervs=https://pvshost:8080,rc=https://rchost:8080;us-south:powervs=https://pvshost-us:8080",
			expectedOutput: []ServiceEndpoint{
				{
					ID:     "vpc",
					URL:    "https://vpchost:8080",
					Region: "eu-gb",
				},
				{
					ID:     "powervs",
					URL:    "https://pvshost:8080",
					Region: "lon",
				},
				{
					ID:     "rc",
					URL:    "https://rchost:8080",
					Region: "lon",
				},
				{
					ID:     "powervs",
					URL:    "https://pvshost-us:8080",
					Region: "us-south",
				},
			},
			expectedError: nil,
		},
		{
			name:           "invalid config",
			flagToParse:    "eu-gb=localhost",
			expectedOutput: nil,
			expectedError:  errServiceEndpointRegion,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := ParseServiceEndpointFlag(tc.flagToParse)
			require.ErrorIs(t, err, tc.expectedError)
			require.ElementsMatch(t, out, tc.expectedOutput)
		})
	}
}
