package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/function61/gokit/app/aws/lambdautils"
	"github.com/function61/gokit/app/dynversion"
	"github.com/function61/gokit/os/osutil"
	"github.com/function61/html2pdf/pkg/h2ptypes"
	"github.com/function61/html2pdf/pkg/html2pdfclient"
	"github.com/spf13/cobra"
)

func main() {
	if lambdautils.InLambda() {
		lambda.StartHandler(lambdautils.NewLambdaHttpHandlerAdapter(newServerHandler()))
		return
	}

	app := &cobra.Command{
		Use:     os.Args[0],
		Short:   "HTML2PDF",
		Version: dynversion.Version,
	}

	app.AddCommand(serverEntry())
	app.AddCommand(clientEntry("client-fn61", html2pdfclient.Function61))
	app.AddCommand(clientEntry("client-localhost", html2pdfclient.Localhost))

	osutil.ExitIfError(app.Execute())
}

func newServerHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/render", func(w http.ResponseWriter, r *http.Request) {
		req := &h2ptypes.Request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			http.Error(w, fmt.Sprintf("Request decode: %v", err), http.StatusBadRequest)
			return
		}

		cmdOpts := h2ptypes.OptionsToWkhtml2PdfCmdline(req.Options)

		// - - => use stdin & stdout
		// (they also need to come after options)
		cmdOpts = append(cmdOpts, "-", "-")

		// we could directly stuff stdout to HTTP response BUT if the command fails, then
		// any output to stdout will end up sending headers and then we can't signal failure
		// anymore, so we'll buffer. we could do this with clever code but that is TODO.
		pdfBuffer := &bytes.Buffer{}

		stdErr := &bytes.Buffer{}

		wkhtmltopdf := exec.Command("./wkhtmltopdf", cmdOpts...)
		wkhtmltopdf.Stdin = bytes.NewReader(req.HtmlBase64)
		wkhtmltopdf.Stdout = pdfBuffer
		wkhtmltopdf.Stderr = stdErr

		if err := wkhtmltopdf.Run(); err != nil {
			http.Error(
				w,
				fmt.Sprintf("%v: %s", err.Error(), stdErr.Bytes()),
				http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")

		if _, err := io.Copy(w, pdfBuffer); err != nil {
			log.Printf("failed to write to client: %v", err)
		}
	})

	return mux
}

func clientEntry(use string, baseUrl string) *cobra.Command {
	return &cobra.Command{
		Use:   use + " [html]",
		Short: "Request HTML2PDF operation from server",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			osutil.ExitIfError(client(
				osutil.CancelOnInterruptOrTerminate(nil),
				args[0],
				os.Stdout,
				baseUrl))
		},
	}
}

func client(ctx context.Context, html string, output io.Writer, baseUrl string) error {
	h2p, err := html2pdfclient.New(baseUrl, html2pdfclient.TokenFromEnv)
	if err != nil {
		return err
	}

	pdf, err := h2p.Render(ctx, html, nil)
	if err != nil {
		return err
	}
	defer pdf.Close()

	_, err = io.Copy(output, pdf)
	return err
}
