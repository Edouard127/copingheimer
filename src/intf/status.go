package intf

import (
	"encoding/json"
)

type StatusResponse struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	}
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

func (s *StatusResponse) ReadFrom(b []byte) error {
	if err := json.Unmarshal(b, s); err != nil {
		// Check if description is an array
	}
	return nil
}
