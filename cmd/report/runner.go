package report

import (
	"context"
	"io"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"

	"github.com/giantswarm/resource-police/pkg/cortex"
	"github.com/giantswarm/resource-police/pkg/intersect"
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
	var err error

	var errors []string

	var clusters []string
	{
		cortexService, err := cortex.New(cortex.Config{
			URL:      r.flag.CortexEndpoint,
			UserName: r.flag.CortexUsername,
			Password: r.flag.CortexPassword,
		})
		if err != nil {
			return microerror.Mask(err)
		}

		// Fetch the clusters that exist right now.
		now := time.Now()
		clusters, err = cortexService.QueryClusters(now)
		if err != nil {
			return microerror.Mask(err)
		}

		// Fetch the clusters that existed three hours ago.
		clustersEarlier, err := cortexService.QueryClusters(now.Add(-(time.Hour * 3)))
		if err != nil {
			return microerror.Mask(err)
		}

		// As a result, report the clusters from 3 hours ago that still exist now.
		clusters = intersect.StringSliceSorted(clusters, clustersEarlier)
	}

	report, err := report.Render(clusters, errors)
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

	// fmt.Println("Report preview:")
	// fmt.Println(report)
	// return nil

	err = slackService.SendReport(report)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
