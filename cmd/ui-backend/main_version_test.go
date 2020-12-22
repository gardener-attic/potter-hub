/*
Copyright (c) 2018 Bitnami

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

package main

import (
	"testing"

	"github.com/arschles/assert"
)

func Test_userAgent(t *testing.T) {
	tests := []testUserAgent{
		// Shows default User-Agent unless comment nor version provided
		{
			version:          "",
			userAgentComment: "",
			expectedResult:   "tiller-proxy/devel",
		},
		// Shows just custom version unless comment provided
		{
			version:          "v4.4.4",
			userAgentComment: "",
			expectedResult:   "tiller-proxy/v4.4.4",
		},
		// Shows custom version plus comment if provided
		{
			version:          "v4.4.4",
			userAgentComment: "Kubeapps/v2.3.4",
			expectedResult:   "tiller-proxy/v4.4.4 (Kubeapps/v2.3.4)",
		},
	}

	for _, tt := range tests {
		run(t, tt)
	}
}

type testUserAgent struct {
	version          string
	userAgentComment string
	expectedResult   string
}

func run(t *testing.T, tt testUserAgent) {
	assert.Equal(t, tt.expectedResult, userAgent(tt.userAgentComment, tt.version), "expected user agent")
}
