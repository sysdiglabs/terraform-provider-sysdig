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
	_ = json.Unmarshal(body, &result)

	return result.User
}

type userWrapper struct {
	User User `json:"user"`
}
