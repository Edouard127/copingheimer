package intf

import (
	"encoding/base64"
	"encoding/json"
)

type StatusResponse struct {
	IP      string `json:"ip"`
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int      `json:"max"`
		Online int      `json:"online"`
		Sample []Player `json:"sample"`
	} `json:"players"`
	Favicon           string `json:"favicon"`
	PreviewChat       string `json:"preview_chat"`
	EnforceSecureChat bool   `json:"enforce_secure_chat"`
}

type Player struct {
	Name string `json:"name"`
	ID   string `json:"id"`
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
