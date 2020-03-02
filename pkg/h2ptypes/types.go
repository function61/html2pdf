package h2ptypes

type Request struct {
	HtmlBase64 []byte   `json:"html_base64"`
	Options    *Options `json:"options"`
}

type Options struct {
	MarginTop    string `json:"marginTop,omitempty"`
	MarginRight  string `json:"marginRight,omitempty"`
	MarginBottom string `json:"marginBottom,omitempty"`
	MarginLeft   string `json:"marginLeft,omitempty"`
	Zoom         string `json:"zoom,omitempty"`
}

func OptionsToWkhtml2PdfCmdline(opts *Options) []string {
	if opts == nil {
		opts = &Options{}
	}

	cmdOpts := []string{}

	cmdOpt := func(key string, value string) {
		if value != "" {
			cmdOpts = append(cmdOpts, key, value)
		}
	}

	cmdOpt("--zoom", opts.Zoom)

	cmdOpt("--margin-top", opts.MarginTop)
	cmdOpt("--margin-right", opts.MarginRight)
	cmdOpt("--margin-bottom", opts.MarginBottom)
	cmdOpt("--margin-left", opts.MarginLeft)

	return cmdOpts
}
