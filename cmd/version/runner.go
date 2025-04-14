package version

import (
	"context"
	"fmt"
	"io"
	"runtime"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"
)

type runner struct {
	flag   *flag
	logger micrologger.Logger
	stdout io.Writer
	stderr io.Writer

	gitCommit string
	source    string
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
	fmt.Printf("Git Commit:     %s\n", r.gitCommit)
	fmt.Printf("Go Version:     %s\n", runtime.Version())
	fmt.Printf("OS / Arch:      %s / %s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Source:         %s\n", r.source)

	return nil
}
