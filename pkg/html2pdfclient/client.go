package html2pdfclient

import (
	"context"
	"fmt"
	"io"

	"github.com/function61/gokit/net/http/ezhttp"
	"github.com/function61/gokit/os/osutil"
	"github.com/function61/html2pdf/pkg/h2ptypes"
)

const (
	Function61 = "https://function61.com/api/html2pdf"
	Localhost  = "http://localhost"
)

type TokenFn func() (string, error)

func TokenFromEnv() (string, error) {
	return osutil.GetenvRequired("HTML2PDF_TOKEN")
}

func NoToken() (string, error) {
	return "", nil
}

type Client struct {
	baseUrl     string
	bearerToken string
}

func New(baseUrl string, getToken TokenFn) (*Client, error) {
	bearerToken, err := getToken()
	if err != nil {
		return nil, fmt.Errorf("getToken: %w", err)
	}

	return &Client{baseUrl, bearerToken}, nil
}

// returns PDF bytes
func (c *Client) Render(
	ctx context.Context,
	html string,
	options *h2ptypes.Options,
) (io.ReadCloser, error) {
	req := &h2ptypes.Request{
		HtmlBase64: []byte(html),
		Options:    options,
	}

	resp, err := ezhttp.Post(
		ctx,
		c.baseUrl+"/render",
		ezhttp.AuthBearer(c.bearerToken),
		ezhttp.Header("Accept", "application/pdf"), // WTF, API gateway returns base64 unless this is set
		ezhttp.SendJson(&req))
	if err != nil {
		return nil, fmt.Errorf("Html2Pdf: %w", err)
	}

	return resp.Body, nil
}
