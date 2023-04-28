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
	NewUsername string
	Id          string
}

type AuthReq struct {
	Password string
	UserName string
}

type RestaurantData struct {
	Name              string
	ImageUrl          string
	OpenTime          string
	CloseTime         string
	Phone             string
	DeliveryCost      float32
	IsDeliveryFixCost bool
	MinDeliveryTime   uint
	MaxDeliveryTime   uint
	DeliveryRadius    float64
}
