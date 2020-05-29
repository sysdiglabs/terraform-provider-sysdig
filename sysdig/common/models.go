package common

import (
	"bytes"
	"encoding/json"
	"io"
)

// -------- User --------
type User struct {
	ID         int    `json:"id,omitempty"`
	Version    int    `json:"version,omitempty"`
	SystemRole string `json:"systemRole,omitempty"`
	Email      string `json:"username"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
}

func (u *User) ToJSON() io.Reader {
	payload, _ := json.Marshal(*u)
	return bytes.NewBuffer(payload)
}

func UserFromJSON(body []byte) User {
	var result userWrapper
	json.Unmarshal(body, &result)

	return result.User
}

type userWrapper struct {
	User User `json:"user"`
}

// -------- Team --------
type Team struct {
	ID                  int         `json:"id,omitempty"`
	Version             int         `json:"version,omitempty"`
	Theme               string      `json:"theme"`
	Name                string      `json:"name"`
	Description         string      `json:"description"`
	ScopeBy             string      `json:"show"`
	Filter              string      `json:"filter"`
	CanUseSysdigCapture bool        `json:"canUseSysdigCapture"`
	CanUseAwsMetrics bool        `json:"canUseAwsMetrics"`
	CanUseCustomEvents bool        `json:"canUseCustomEvents"`
	CanUseBeaconMetrics bool        `json:"canUseBeaconMetrics"`
	UserRoles           []UserRoles `json:"userRoles,omitempty"`
	DefaultTeam         bool        `json:"default"`
	Products            []string    `json:"products"`
}

type UserRoles struct {
	UserId int    `json:"userId"`
	Email  string `json:"userName",omitempty`
	Role   string `json:"role"`
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
