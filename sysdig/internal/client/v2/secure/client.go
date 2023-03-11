package secure

import "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/common"

type client struct {
	*common.Client
}

type Secure interface {
	common.Common
}

var (
	WithURL          = common.WithURL
	WithToken        = common.WithToken
	WithInsecure     = common.WithInsecure
	WithExtraHeaders = common.WithExtraHeaders
)

func New(opt ...common.ClientOption) Secure {
	return &client{
		common.NewClient(opt...),
	}
}
