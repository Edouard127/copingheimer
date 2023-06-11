package core

import (
	pk "github.com/Tnze/go-mc/net/packet"
	"io"
)

type StatusResponse struct {
	IP      string `json:"ip"`
	Version struct {
		Name     string `json:"name"`
		Protocol int32  `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int32    `json:"max"`
		Online int32    `json:"online"`
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

func (p Player) WriteTo(w io.Writer) (int64, error) {
	return pk.Tuple{
		pk.String(p.Name),
		pk.String(p.ID),
	}.WriteTo(w)
}

func (p *Player) ReadFrom(r io.Reader) (int64, error) {
	return pk.Tuple{
		(*pk.String)(&p.Name),
		(*pk.String)(&p.ID),
	}.ReadFrom(r)
}

func (s StatusResponse) WriteTo(w io.Writer) (int64, error) {
	n, err := pk.Tuple{
		pk.String(s.IP),
		pk.Tuple{
			pk.String(s.Version.Name),
			pk.Int(s.Version.Protocol),
		},
		pk.Tuple{
			pk.Int(s.Players.Max),
			pk.Int(s.Players.Online),
			pk.Array(s.Players.Sample),
		},
		pk.String(s.Favicon),
		pk.String(s.PreviewChat),
		pk.Boolean(s.EnforceSecureChat),
	}.WriteTo(w)

	return n, err
}

func (s *StatusResponse) ReadFrom(r io.Reader) (int64, error) {
	n1, err := pk.Tuple{
		(*pk.String)(&s.IP),
		pk.Tuple{
			(*pk.String)(&s.Version.Name),
			(*pk.Int)(&s.Version.Protocol),
		},
		pk.Tuple{
			(*pk.Int)(&s.Players.Max),
			(*pk.Int)(&s.Players.Online),
		},
	}.ReadFrom(r)

	s.Players.Sample = make([]Player, s.Players.Online)
	n2, err := pk.Array(&s.Players.Sample).ReadFrom(r)
	n3, err := pk.Tuple{
		(*pk.String)(&s.Favicon),
		(*pk.String)(&s.PreviewChat),
		(*pk.Boolean)(&s.EnforceSecureChat),
	}.ReadFrom(r)

	return n1 + n2 + n3, err
}
