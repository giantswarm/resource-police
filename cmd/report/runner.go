package report

import (
	"context"
	"io"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"

	"github.com/giantswarm/resource-police/pkg/installation"
	"github.com/giantswarm/resource-police/pkg/slack"
)

type runner struct {
	flag   *flag
	logger micrologger.Logger
	stdout io.Writer
	stderr io.Writer
}

func (r *runner) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	err := r.flag.Validate()
	if err != nil {
		return microerror.Mask(err)
	}

	err = r.run(ctx, cmd, args)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *runner) run(ctx context.Context, cmd *cobra.Command, args []string) error {
	var err error

	var testInstallations []installation.Installation
	{
		c := installation.Config{
			Logger: r.logger,

			InstallationsConfigFile: r.flag.InstallationsConfigFile,
		}

		testInstallations, err = installation.New(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	var collectedInstallations []installation.Installation
	{
		for _, i := range testInstallations {
			clusters, err := installation.ListClusters(i)
			if err != nil {
				return microerror.Mask(err)
			}

			if len(clusters) > 0 {
				newInstallation := installation.Installation{
					Name:     i.Name,
					Clusters: clusters,
				}
				collectedInstallations = append(collectedInstallations, newInstallation)
			}
		}
	}

	report, err := installation.RenderReport(collectedInstallations)
	if err != nil {
		return microerror.Mask(err)
	}

	var slackService *slack.Slack
	{
		c := slack.Config{
			Logger: r.logger,

			WebhookEndpoint: r.flag.SlackWebhookEndpoint,
		}

		slackService, err = slack.New(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	err = slackService.SendReport(report)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
