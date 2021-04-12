package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

type Config struct {
	Logger micrologger.Logger

	WebhookEndpoint string
}

type Slack struct {
	logger micrologger.Logger

	webhookEndpoint string
}

func New(config Config) (*Slack, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.WebhookEndpoint == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.WebhookEndpoint must not be empty", config)
	}

	slack := &Slack{
		logger:          config.Logger,
		webhookEndpoint: config.WebhookEndpoint,
	}

	return slack, nil
}

func (s *Slack) SendReport(report string) error {
	fmt.Println("Sending report to slack...")

	requestBody, err := json.Marshal(map[string]string{
		"text": report,
	})
	if err != nil {
		return microerror.Mask(err)
	}

	_, err = http.Post(s.webhookEndpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return microerror.Mask(err)
	}

	fmt.Println("Report has been sent.")

	return nil
}
