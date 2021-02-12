package main

import (
	"context"
	"log"
	"net/http"

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

	return httputils.CancelableServer(ctx, srv, func() error { return srv.ListenAndServe() })
}
