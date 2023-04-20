package types

import "github.com/nicodeheza/peersEat/models"

type PeerPresentationBody struct {
	NewPeer models.Peer
	SendTo  []string
}

type SendAllPeerQuery struct {
	Excludes []string `query:"excludes"`
}

type ApiGeoCord struct{}

type GetCordsResponse struct {
	Place_id     int
	Licence      string
	Osm_type     string
	Osm_id       int
	Boundingbox  []string
	Lat          string
	Lon          string
	Display_name string
	Class        string
	Type         string
	Importance   float64
}

type PeerHaveRestaurantResp struct {
	Resp bool
	Err  error
}

type UpdateRestaurantPassword struct {
	NewPassword string
	Id          string
}

type AuthReq struct {
	Password string
	UserName string
}
