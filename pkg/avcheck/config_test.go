package avcheck

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/arschles/assert"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      Configuration
		producesErr bool
		errMsg      string
	}{
		{
			name: "valid",
			config: Configuration{
				PathPrefix:                   "/potter-health",
				ChartsAvailableCheckInterval: minChartsAvailableCheckInterval,
			},
			producesErr: false,
			errMsg:      "",
		},
		{
			name: "ChartsAvailableCheckInterval too small",
			config: Configuration{
				PathPrefix:                   "/potter-health",
				ChartsAvailableCheckInterval: 5 * time.Millisecond,
			},
			producesErr: true,
			errMsg:      fmt.Sprintf("chartsAvailableCheckInterval must be greater than %s", minChartsAvailableCheckInterval),
		},
		{
			name: "PathPrefix is empty",
			config: Configuration{
				PathPrefix:                   "",
				ChartsAvailableCheckInterval: 5 * time.Millisecond,
			},
			producesErr: true,
			errMsg:      "pathPrefix must not be empty",
		},
		{
			name: "PathPrefix not starting with a /",
			config: Configuration{
				PathPrefix:                   "potter-health",
				ChartsAvailableCheckInterval: 5 * time.Millisecond,
			},
			producesErr: true,
			errMsg:      "pathPrefix must start with a /",
		},
		{
			name: "PathPrefix ends with a /",
			config: Configuration{
				PathPrefix:                   "/potter-health/",
				ChartsAvailableCheckInterval: 5 * time.Millisecond,
			},
			producesErr: true,
			errMsg:      "pathPrefix must not end with a /",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.producesErr {
				assert.NotNil(t, err, "validation error")
				assert.Equal(t, err.Error(), tt.errMsg, "error message")
			} else {
				assert.Nil(t, err, "validation error")
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name           string
		configJSON     string
		producesErr    bool
		expectedConfig Configuration
	}{
		{
			name: "valid",
			configJSON: `{
				"pathPrefix": "/potter-health",
				"chartsAvailableCheckInterval": "15m"
			}`,
			producesErr: false,
			expectedConfig: Configuration{
				PathPrefix:                   "/potter-health",
				ChartsAvailableCheckInterval: 15 * time.Minute,
			},
		},
		{
			name: "invalid pathPrefix datatype",
			configJSON: `{
				"pathPrefix": 14,
				"chartsAvailableCheckInterval": "15m"
			}`,
			producesErr: true,
		},
		{
			name: "invalid chartsAvailableCheckInterval datatype",
			configJSON: `{
				"pathPrefix": "/potter-health",
				"chartsAvailableCheckInterval": 15
			}`,
			producesErr: true,
		},
		{
			name: "invalid chartsAvailableCheckInterval value",
			configJSON: `{
				"pathPrefix": "/potter-health",
				"chartsAvailableCheckInterval": "thisIsNoTime"
			}`,
			producesErr: true,
		},
		{
			name: "invalid JSON structure",
			configJSON: `{
				"pathPrefix": "/potter-health",
				"chartsAvailableCheckInterval": "15m",
				{}
			}`,
			producesErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var config Configuration
			err := json.Unmarshal([]byte(tt.configJSON), &config)
			if tt.producesErr {
				assert.NotNil(t, err, "unmarshaling error")
			} else {
				assert.Nil(t, err, "unmarshaling error")
				assert.Equal(t, config, tt.expectedConfig, "config")
			}
		})
	}
}
