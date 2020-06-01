package monitor

import (
	"bytes"
	"encoding/json"
	"io"
)

// -------- Alert --------
type CustomNotification struct {
	TitleTemplate  string `json:"titleTemplate"`
	UseNewTemplate bool   `json:"useNewTemplate"`
}

type SysdigCapture struct {
	Name       string      `json:"name"`
	Filters    string      `json:"filters,omitempty"`
	Duration   int         `json:"duration"`
	Type       string      `json:"type,omitempty"`
	BucketName string      `json:"bucketName"`
	Folder     string      `json:"folder,omitempty"`
	Enabled    bool        `json:"enabled"`
	StorageID  interface{} `json:"storageId,omitempty"`
}
type SegmentCondition struct {
	Type string `json:"type"`
}

type Criteria struct {
	Text   string `json:"text"`
	Source string `json:"source"`
}

type Monitor struct {
	Metric       string  `json:"metric"`
	StdDevFactor float64 `json:"stdDevFactor"`
}

type alertWrapper struct {
	Alert Alert `json:"alert"`
}

type Alert struct {
	ID                     int                 `json:"id,omitempty"`
	Version                int                 `json:"version,omitempty"`
	Type                   string              `json:"type"` // computed MANUAL
	Name                   string              `json:"name"`
	Description            string              `json:"description"`
	Enabled                bool                `json:"enabled"`
	NotificationChannelIds []int               `json:"notificationChannelIds"`
	Filter                 string              `json:"filter"`
	Severity               int                 `json:"severity"` // 6 == INFO, 4 == LOW, 2 == MEDIUM, 0 == HIGH // NOT USED
	Timespan               int                 `json:"timespan"` // computed 600000000
	CustomNotification     *CustomNotification `json:"customNotification"`
	TeamID                 int                 `json:"teamId,omitempty"` // computed
	AutoCreated            bool                `json:"autoCreated"`
	SysdigCapture          *SysdigCapture      `json:"sysdigCapture"`
	RateOfChange           bool                `json:"rateOfChange,omitempty"`
	ReNotifyMinutes        int                 `json:"reNotifyMinutes"`
	ReNotify               bool                `json:"reNotify"`
	Valid                  bool                `json:"valid"`
	SeverityLabel          string              `json:"severityLabel,omitempty"` // MEDIUM == MEDIUM, LOW == LOW, NONE == INFO, HIGH == HIGH
	SegmentBy              []string            `json:"segmentBy"`
	SegmentCondition       *SegmentCondition   `json:"segmentCondition"`
	Criteria               *Criteria           `json:"criteria,omitempty"`
	Monitor                []*Monitor          `json:"monitor,omitempty"`
	Condition              string              `json:"condition"`
	SeverityLevel          int                 `json:"severityLevel,omitempty"` // 0 == MEDIUM, 2 == LOW, 4 == INFO, 6 == HIGH
}

func (a *Alert) ToJSON() io.Reader {
	payload, _ := json.Marshal(alertWrapper{Alert: *a})
	return bytes.NewBuffer(payload)
}

func AlertFromJSON(body []byte) Alert {
	var result alertWrapper
	json.Unmarshal(body, &result)

	return result.Alert
}

// -------- Team --------
type Team struct {
	UserRoles           []UserRoles `json:"userRoles"`
	Description         string      `json:"description"`
	Name                string      `json:"name"`
	ID                  int         `json:"id,omitempty"`
	Default             bool        `json:"default"`
	Version             int         `json:"version,omitempty"`
	Origin              string      `json:"origin,omitempty"`
	LastUpdated         int64       `json:"lastUpdated,omitempty"`
	EntryPoint          EntryPoint  `json:"entryPoint"`
	Theme               string      `json:"theme"`
	CustomerID          int         `json:"customerId"`
	DateCreated         int64       `json:"dateCreated"`
	Products            []string    `json:"products,omitempty"`
	Show                string      `json:"show"`
	Immutable           bool        `json:"immutable"`
	CanUseSysdigCapture bool        `json:"canUseSysdigCapture"`
	CanUseCustomEvents  bool        `json:"canUseCustomEvents"`
	CanUseAwsMetrics    bool        `json:"canUseAwsMetrics"`
	CanUseBeaconMetrics bool        `json:"canUseBeaconMetrics"`
	UserCount           int         `json:"userCount"`
	Filter              string      `json:"filter,omitempty"`
	DefaultTeam         bool        `json:"default,omitempty"`
}

type UserRoles struct {
	UserId int    `json:"userId"`
	Email  string `json:"userName,omitempty"`
	Role   string `json:"role"`
}

type EntryPoint struct {
	Module    string `json:"module"`
	Selection string `json:"selection,omitempty"`
}

func (t *Team) ToJSON() io.Reader {
	payload, _ := json.Marshal(*t)
	return bytes.NewBuffer(payload)
}

func TeamFromJSON(body []byte) Team {
	var result teamWrapper
	json.Unmarshal(body, &result)

	return result.Team
}

type teamWrapper struct {
	Team Team `json:"team"`
}

// -------- UsersList --------
type UsersList struct {
	ID    int    `json:"id"`
	Email string `json:"username"`
}

func UsersListFromJSON(body []byte) []UsersList {
	var result usersListWrapper
	json.Unmarshal(body, &result)

	return result.UsersList
}

type usersListWrapper struct {
	UsersList []UsersList `json:"users"`
}
