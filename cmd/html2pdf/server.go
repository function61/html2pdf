package main

import (
	"context"
	"github.com/function61/gokit/httputils"
	"github.com/function61/gokit/logex"
	"github.com/function61/gokit/ossignal"
	"github.com/function61/gokit/taskrunner"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

func serverEntry() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start server (also good for dev/testing)",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger := logex.StandardLogger()

			exitIfError(runServer(
				ossignal.InterruptOrTerminateBackgroundCtx(logger),
				logger))
		},
	}
}

func runServer(ctx context.Context, logger *log.Logger) error {
	srv := &http.Server{
		Addr:    ":80",
		Handler: newServerHandler(),
	}

	tasks := taskrunner.New(ctx, logger)

	tasks.Start("listener "+srv.Addr, func(_ context.Context, _ string) error {
		return httputils.RemoveGracefulServerClosedError(srv.ListenAndServe())
	})

	tasks.Start("listenershutdowner", httputils.ServerShutdownTask(srv))

	return tasks.Wait()
}
