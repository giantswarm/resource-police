package report

import (
	"context"
	"fmt"
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
	var errors []error

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
		clustersNow, err := cortexService.QueryClusters(now)
		if err != nil {
			errors = append(errors, err)
		}

		// Fetch the clusters that existed three hours ago.
		// That's the age threshold we use for reporting a test cluster.
		threeHoursAgo := now.Add(-(time.Hour * 3))
		clustersEarlier, err := cortexService.QueryClusters(threeHoursAgo)
		if err != nil {
			errors = append(errors, err)
		}

		// Create intersection of both queries.
		clusters = intersect.StringSlice(clustersNow, clustersEarlier)
	}

	report, err := report.Render(clusters, errors)
	if err != nil {
		return microerror.Mask(err)
	}

	fmt.Println("Report preview:")
	fmt.Println(report)
	return nil

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
