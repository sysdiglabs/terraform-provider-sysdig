package secure

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

type SysdigSecureClient interface {
	CreatePolicy(Policy) (Policy, error)
	DeletePolicy(int) error
	UpdatePolicy(Policy) (Policy, error)
	GetPolicyById(int) (Policy, error)

	CreateRule(Rule) (Rule, error)
	GetRuleByID(int) (Rule, error)
	UpdateRule(Rule) (Rule, error)
	DeleteRule(int) error

	CreateNotificationChannel(NotificationChannel) (NotificationChannel, error)
	GetNotificationChannelById(int) (NotificationChannel, error)
	GetNotificationChannelByName(string) (NotificationChannel, error)
	DeleteNotificationChannel(int) error
	UpdateNotificationChannel(NotificationChannel) (NotificationChannel, error)

	CreateUser(User) (User, error)
	GetUserById(int) (User, error)
	DeleteUser(int) error
	UpdateUser(User) (User, error)

	CreateTeam(Team) (Team, error)
	GetTeamById(int) (Team, error)
	DeleteTeam(int) error
	UpdateTeam(Team) (Team, error)

	CreateList(List) (List, error)
	GetListById(int) (List, error)
	DeleteList(int) error
	UpdateList(List) (List, error)
}

func NewSysdigSecureClient(sysdigSecureAPIToken string, url string) SysdigSecureClient {
	return &sysdigSecureClient{
		SysdigSecureAPIToken: sysdigSecureAPIToken,
		URL:                  url,
		httpClient:           http.DefaultClient,
	}
}

type sysdigSecureClient struct {
	SysdigSecureAPIToken string
	URL                  string
	httpClient           *http.Client
}

func (client *sysdigSecureClient) doSysdigSecureRequest(method string, url string, payload io.Reader) (*http.Response, error) {
	request, _ := http.NewRequest(method, url, payload)
	request.Header.Set("Authorization", "Bearer "+client.SysdigSecureAPIToken)
	request.Header.Set("Content-Type", "application/json")

	out, _ := httputil.DumpRequestOut(request, true)
	log.Printf("[DEBUG] %s", string(out))
	response, error := client.httpClient.Do(request)

	out, _ = httputil.DumpResponse(response, true)
	log.Printf("[DEBUG] %s", string(out))
	return response, error
}
