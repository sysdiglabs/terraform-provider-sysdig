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
	ID         int      `json:"id,omitempty"`
	GroupName  string   `json:"groupName,omitempty"`
	Role       string   `json:"role,omitempty"`
	SystemRole string   `json:"systemRole,omitempty"`
	TeamMap    *TeamMap `json:"teamMap,omitempty"`
}

func (gm *GroupMapping) ToJSON() (io.Reader, error) {
	payload, err := json.Marshal(*gm)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(payload), nil
}
