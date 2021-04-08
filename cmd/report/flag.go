package report

import (
	"os"

	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"

	"github.com/giantswarm/resource-police/internal/env"
)

const (
	flagSlackWebhookEndpoint = "slack.webhook.endpoint"

	flagCortexEndpoint = "cortex.endpoint.url"
	flagCortexUsername = "cortex.username"
	flagCortexPassword = "cortex.password"
)

type flag struct {
	SlackWebhookEndpoint string

	CortexEndpoint string
	CortexUsername string
	CortexPassword string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.SlackWebhookEndpoint, flagSlackWebhookEndpoint, os.Getenv(env.SlackWebhookEndpoint), "Slack Webhook endpoint for posting messages into channel")
	cmd.Flags().StringVar(&f.CortexEndpoint, flagCortexEndpoint, "https://prometheus-us-central1.grafana.net/api/prom", "Cortex endpoint URL")
	cmd.Flags().StringVar(&f.CortexUsername, flagCortexUsername, os.Getenv(env.CortexUserName), "Cortex user ID")
	cmd.Flags().StringVar(&f.CortexPassword, flagCortexPassword, os.Getenv(env.CortexPassword), "Cortex API token")
}

func (f *flag) Validate() error {
	if f.SlackWebhookEndpoint == "" {
		return microerror.Maskf(invalidFlagError, "--%s or %s environment variable must not be empty", flagSlackWebhookEndpoint, env.SlackWebhookEndpoint)
	}
	if f.CortexEndpoint == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagCortexEndpoint)
	}
	if f.CortexUsername == "" {
		return microerror.Maskf(invalidFlagError, "--%s or %s environment variable must not be empty", flagCortexUsername, env.CortexUserName)
	}
	if f.CortexPassword == "" {
		return microerror.Maskf(invalidFlagError, "--%s or %s environment variable must not be empty", flagCortexPassword, env.CortexPassword)
	}

	return nil
}
