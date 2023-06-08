package v2

type config struct {
	url            string
	token          string
	insecure       bool
	extraHeaders   map[string]string
	ibmInstanceID  string
	ibmAPIKey      string
	ibmIamURL      string
	sysdigTeamName string
	sysdigTeamID   *int
	product        string
}

type Product string

const (
	MonitorProduct Product = "SDC"
	SecureProduct  Product = "SDS"
)

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

func WithIBMInstanceID(instanceID string) ClientOption {
	return func(c *config) {
		c.ibmInstanceID = instanceID
	}
}

func WithIBMAPIKey(key string) ClientOption {
	return func(c *config) {
		c.ibmAPIKey = key
	}
}

func WithIBMIamURL(url string) ClientOption {
	return func(c *config) {
		c.ibmIamURL = url
	}
}

func WithSysdigTeamID(teamID *int) ClientOption {
	return func(c *config) {
		c.sysdigTeamID = teamID
	}
}

func WithSysdigTeamName(teamName string) ClientOption {
	return func(c *config) {
		c.sysdigTeamName = teamName
	}
}

func WithMonitorProduct() ClientOption {
	return func(c *config) {
		c.product = string(MonitorProduct)
	}
}

func WithSecureProduct() ClientOption {
	return func(c *config) {
		c.product = string(SecureProduct)
	}
}

func configure(opts ...ClientOption) *config {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
