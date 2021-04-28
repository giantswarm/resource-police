package report

import (
	"context"
	"fmt"
	"io"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"

	"github.com/giantswarm/resource-police/pkg/cortex"
	"github.com/giantswarm/resource-police/pkg/report"
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
	var errors []error

	var clusters []cortex.Cluster
	{
		cortexService, err := cortex.New(cortex.Config{
			URL:      r.flag.CortexEndpoint,
			UserName: r.flag.CortexUsername,
			Password: r.flag.CortexPassword,
		})
		if err != nil {
			return microerror.Mask(err)
		}

		// Fetch the clusters info
		clusters, err = cortexService.QueryClusters()
		if err != nil {
			errors = append(errors, err)
		}
	}

	report, err := report.Render(clusters, errors)
	if err != nil {
		return microerror.Mask(err)
	}

	if r.flag.DryRun {
		fmt.Printf("\nReport:\n\n")
		fmt.Println(report)
		return nil
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
