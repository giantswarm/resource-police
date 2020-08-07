package report

import (
	"context"
	"fmt"
	"io"
	"sort"

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

	var errors []string

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

	var clustersToDelete []*installation.Cluster
	{
		for _, i := range testInstallations {
			clusters, err := installation.ListClusters(i)
			if err == nil {
				clustersToDelete = append(clustersToDelete, clusters...)
			} else {
				// collect errors, but continue
				errors = append(errors, fmt.Sprintf("Could not list clusters in installation `%s`: `%s`", i.Name, err))
			}
		}
	}

	sort.Slice(clustersToDelete, func(i, j int) bool {
		return clustersToDelete[i].Age.Milliseconds() > clustersToDelete[j].Age.Milliseconds()
	})

	report, err := installation.RenderReport(clustersToDelete, errors)
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
