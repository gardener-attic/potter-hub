package avcheck

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const minChartsAvailableCheckInterval = time.Second * 10

type Configuration struct {
	PathPrefix                   string
	ChartsAvailableCheckInterval time.Duration
}

func (c *Configuration) Validate() error {
	if c.PathPrefix == "" {
		return errors.New("pathPrefix must not be empty")
	}

	if !strings.HasPrefix(c.PathPrefix, "/") {
		return errors.New("pathPrefix must start with a /")
	}

	if strings.HasSuffix(c.PathPrefix, "/") {
		return errors.New("pathPrefix must not end with a /")
	}

	if c.ChartsAvailableCheckInterval < minChartsAvailableCheckInterval {
		return errors.Errorf("chartsAvailableCheckInterval must be greater than %s", minChartsAvailableCheckInterval)
	}

	return nil
}

func (c *Configuration) UnmarshalJSON(data []byte) error {
	var tmp struct {
		PathPrefix                   string `json:"pathPrefix"`
		ChartsAvailableCheckInterval string `json:"chartsAvailableCheckInterval"`
	}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	c.PathPrefix = tmp.PathPrefix

	c.ChartsAvailableCheckInterval, err = time.ParseDuration(tmp.ChartsAvailableCheckInterval)
	if err != nil {
		return err
	}

	return nil
}
