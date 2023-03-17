package v2

type config struct {
	url          string
	token        string
	insecure     bool
	extraHeaders map[string]string
}

type ClientOption func(c *config)

func WithURL(url string) ClientOption {
	return func(c *config) {
		c.url = url
	}
}

func WithToken(token string) ClientOption {
	return func(c *config) {
		c.token = token
	}
}

func WithInsecure(insecure bool) ClientOption {
	return func(c *config) {
		c.insecure = insecure
	}
}

func WithExtraHeaders(headers map[string]string) ClientOption {
	return func(c *config) {
		c.extraHeaders = headers
	}
}

func configure(opts ...ClientOption) *config {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
