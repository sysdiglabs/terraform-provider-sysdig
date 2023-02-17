package monitor

import (
	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/common"
)

type client struct {
	*common.Client
}
type Client interface {
	SetExtraHeaders(extraHeaders map[string]string)
	common.GroupMapper
}

func NewClient(token string, url string, insecure bool) Client {
	return &client{
		Client: common.NewClient(token, url, insecure),
	}
}
