package common

type Team struct {
	UserRoles           []UserRoles `json:"userRoles,omitempty"`
	Description         string      `json:"description,omitempty"`
	Name                string      `json:"name,omitempty"`
	ID                  int         `json:"id,omitempty"`
	Version             int         `json:"version,omitempty"`
	Origin              string      `json:"origin,omitempty"`
	LastUpdated         int64       `json:"lastUpdated,omitempty"`
	EntryPoint          *EntryPoint `json:"entryPoint,omitempty"`
	Theme               string      `json:"theme,omitempty"`
	CustomerID          int         `json:"customerId,omitempty"`
	DateCreated         int64       `json:"dateCreated,omitempty"`
	Products            []string    `json:"products,omitempty"`
	Show                string      `json:"show,omitempty"`
	Immutable           bool        `json:"immutable,omitempty"`
	CanUseSysdigCapture bool        `json:"canUseSysdigCapture,omitempty"`
	CanUseCustomEvents  bool        `json:"canUseCustomEvents,omitempty"`
	CanUseAwsMetrics    bool        `json:"canUseAwsMetrics,omitempty"`
	CanUseBeaconMetrics bool        `json:"canUseBeaconMetrics,omitempty"`
	UserCount           int         `json:"userCount,omitempty"`
	Filter              string      `json:"filter,omitempty"`
	DefaultTeam         bool        `json:"default,omitempty"`
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

type UsersList struct {
	ID    int    `json:"id"`
	Email string `json:"username"`
}

type usersListWrapper struct {
	UsersList []UsersList `json:"users"`
}
