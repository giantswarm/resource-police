package report

import (
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"

	"github.com/giantswarm/resource-police/internal/env"
)

const (
	flagInstallationsConfigFile = "installations.config.file"
	flagSlackWebhookEndpoint    = "slack.webhook.endpoint"
)

type flag struct {
	InstallationsConfigFile string
	SlackWebhookEndpoint    string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.InstallationsConfigFile, flagInstallationsConfigFile, "", "Installations configuration file.")
	cmd.Flags().StringVar(&f.SlackWebhookEndpoint, flagSlackWebhookEndpoint, "", "Slack Webhook endpoint for posting messages into channel.")
}

func (f *flag) Validate() error {

	if f.InstallationsConfigFile == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagInstallationsConfigFile)
	}
	if f.SlackWebhookEndpoint == "" {
		return microerror.Maskf(invalidFlagError, "--%s or %s environment variable must not be empty", flagSlackWebhookEndpoint, env.SlackWebhookEndpoint)
	}

	return nil
}
