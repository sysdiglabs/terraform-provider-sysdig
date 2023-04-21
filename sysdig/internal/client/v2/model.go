package v2

type User struct {
	ID          int `json:"id"`
	CurrentTeam int `json:"currentTeam"`
}

type userWrapper struct {
	User User `json:"user"`
}

type Team struct {
	UserRoles           []UserRoles       `json:"userRoles,omitempty"`
	Description         string            `json:"description"`
	Name                string            `json:"name"`
	ID                  int               `json:"id,omitempty"`
	Version             int               `json:"version,omitempty"`
	Origin              string            `json:"origin,omitempty"`
	LastUpdated         int64             `json:"lastUpdated,omitempty"`
	EntryPoint          *EntryPoint       `json:"entryPoint,omitempty"`
	Theme               string            `json:"theme"`
	CustomerID          int               `json:"customerId,omitempty"`
	DateCreated         int64             `json:"dateCreated,omitempty"`
	Products            []string          `json:"products,omitempty"`
	Show                string            `json:"show,omitempty"`
	Immutable           bool              `json:"immutable,omitempty"`
	CanUseSysdigCapture *bool             `json:"canUseSysdigCapture,omitempty"`
	CanUseCustomEvents  *bool             `json:"canUseCustomEvents,omitempty"`
	CanUseAwsMetrics    *bool             `json:"canUseAwsMetrics,omitempty"`
	CanUseBeaconMetrics *bool             `json:"canUseBeaconMetrics,omitempty"`
	UserCount           int               `json:"userCount,omitempty"`
	Filter              string            `json:"filter,omitempty"`
	NamespaceFilters    *NamespaceFilters `json:"namespaceFilters,omitempty"`
	DefaultTeam         bool              `json:"default,omitempty"`
}

type NamespaceFilters struct {
	IBMPlatformMetrics *string `json:"ibmPlatformMetrics"`
}

type UserRoles struct {
	UserId int    `json:"userId"`
	Email  string `json:"userName,omitempty"`
	Role   string `json:"role"`
	Admin  bool   `json:"admin,omitempty"`
}

type EntryPoint struct {
	Module    string `json:"module"`
	Selection string `json:"selection,omitempty"`
}

type teamWrapper struct {
	Team Team `json:"team"`
}

type User struct {
	ID         int    `json:"id,omitempty"`
	Version    int    `json:"version,omitempty"`
	SystemRole string `json:"systemRole,omitempty"`
	Email      string `json:"username"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
}

type userWrapper struct {
	User User `json:"user"`
}

type usersWrapper struct {
	Users []User `json:"users"`
}

type NotificationChannelOptions struct {
	EmailRecipients   []string               `json:"emailRecipients,omitempty"`   // Type: email
	SnsTopicARNs      []string               `json:"snsTopicARNs,omitempty"`      // Type: SNS
	APIKey            string                 `json:"apiKey,omitempty"`            // Type: VictorOps
	RoutingKey        string                 `json:"routingKey,omitempty"`        // Type: VictorOps
	Url               string                 `json:"url,omitempty"`               // Type: OpsGenie, Webhook and Slack
	Channel           string                 `json:"channel,omitempty"`           // Type: Slack
	Account           string                 `json:"account,omitempty"`           // Type: PagerDuty
	ServiceKey        string                 `json:"serviceKey,omitempty"`        // Type: PagerDuty
	ServiceName       string                 `json:"serviceName,omitempty"`       // Type: PagerDuty
	AdditionalHeaders map[string]interface{} `json:"additionalHeaders,omitempty"` // Type: Webhook
	Region            string                 `json:"region,omitempty"`            // Type: OpsGenie

	NotifyOnOk           bool `json:"notifyOnOk"`
	NotifyOnResolve      bool `json:"notifyOnResolve"`
	SendTestNotification bool `json:"sendTestNotification"`
}

type NotificationChannel struct {
	ID      int                        `json:"id,omitempty"`
	Version int                        `json:"version,omitempty"`
	Type    string                     `json:"type"`
	Name    string                     `json:"name"`
	Enabled bool                       `json:"enabled"`
	TeamID  *int                       `json:"teamId,omitempty"`
	Options NotificationChannelOptions `json:"options"`
}

type notificationChannelListWrapper struct {
	NotificationChannels []NotificationChannel `json:"notificationChannels"`
}

type notificationChannelWrapper struct {
	NotificationChannel NotificationChannel `json:"notificationChannel"`
}
