package slack

import (
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
	return nil
}
