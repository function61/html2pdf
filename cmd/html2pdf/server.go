package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/function61/gokit/log/logex"
	"github.com/function61/gokit/net/http/httputils"
	"github.com/function61/gokit/os/osutil"
	"github.com/spf13/cobra"
)

func serverEntry() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start server (also good for dev/testing)",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger := logex.StandardLogger()

			osutil.ExitIfError(runServer(
				osutil.CancelOnInterruptOrTerminate(logger),
				logger))
		},
	}
}

func runServer(ctx context.Context, logger *log.Logger) error {
	srv := &http.Server{
		Addr:    ":80",
		Handler: newServerHandler(),
	}

	// this codepath we're on right now is not used in Lambda
	if err := fixFontsWhenNotRunningInLambda(); err != nil {
		return err
	}

	return httputils.CancelableServer(ctx, srv, srv.ListenAndServe)
}

const (
	fontsPathOnLambda = "/var/task/fonts"
)

// on Lambda our fonts are stored under /var/task tree, which might not match local development.
// it's easiest to use a symlink during local dev to mimic Lambda environment.
func fixFontsWhenNotRunningInLambda() error {
	fontsDirExists, err := osutil.Exists(fontsPathOnLambda)
	if err != nil {
		return err
	}

	if !fontsDirExists {
		if err := os.MkdirAll(filepath.Dir(fontsPathOnLambda), 0775); err != nil {
			return err
		}

		if err := os.Symlink("/workspace/fonts", fontsPathOnLambda); err != nil {
			return err
		}
	}

	return nil
}
