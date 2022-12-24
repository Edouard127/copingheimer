package intf

import (
	"encoding/base64"
	"encoding/json"
)

type StatusResponse struct {
	IP      string `json:"ip" bson:"ip"`
	Version struct {
		Name     string `json:"name" bson:"name"`
		Protocol int    `json:"protocol" bson:"protocol"`
	} `json:"version" bson:"version"`
	Players struct {
		Max    int      `json:"max" bson:"max"`
		Online int      `json:"online" bson:"online"`
		Sample []Player `json:"sample" bson:"sample"`
	} `json:"players" bson:"players"`
	Favicon           string `json:"favicon" bson:"favicon"`
	PreviewChat       string `json:"preview_chat" bson:"preview_chat"`
	EnforceSecureChat bool   `json:"enforce_secure_chat" bson:"enforce_secure_chat"`
}

type Player struct {
	Name string `json:"name" bson:"name"`
	ID   string `json:"id" bson:"id"`
}

func (s *StatusResponse) String() string {
	b, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

func (s *StatusResponse) Json(ip string) ([]byte, error) {
	s.IP = ip
	return json.Marshal(s)
}

func (s *StatusResponse) Put(b []byte) error {
	return json.Unmarshal(b, s)
}
