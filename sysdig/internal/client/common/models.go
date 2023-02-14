package common

import (
	"bytes"
	"encoding/json"
	"io"
)

// -------- Group mapping --------
type TeamMap struct {
	AllTeams bool  `json:"allTeams"`
	TeamIDs  []int `json:"teamIds"`
}

type GroupMapping struct {
	ID        int      `json:"id,omitempty"`
	GroupName string   `json:"groupName,omitempty"`
	Role      string   `json:"role,omitempty"`
	TeamMap   *TeamMap `json:"teamMap,omitempty"`
}

func (gm *GroupMapping) ToJSON() (io.Reader, error) {
	payload, err := json.Marshal(*gm)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(payload), nil
}

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
	_ = json.Unmarshal(body, &result)

	return result.User
}

type userWrapper struct {
	User User `json:"user"`
}
