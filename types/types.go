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
	Id                string  `validate:"required"`
	Name              string  `validate:"required"`
	ImageUrl          string  `validate:"required,url"`
	OpenTime          string  `validate:"required,dateTime=3:04PM"`
	CloseTime         string  `validate:"required,dateTime=3:04PM"`
	Phone             string  `validate:"required,e164"`
	DeliveryCost      float32 `validate:"required,gte=0"`
	IsDeliveryFixCost bool    `validate:"required"`
	MinDeliveryTime   uint    `validate:"required,gte=0"`
	MaxDeliveryTime   uint    `validate:"required,gte=0"`
	DeliveryRadius    float64 `validate:"required,gte=0"`
}

type Event struct {
	Name    string      `validate:"required"`
	Payload interface{} `validate:"required"`
	SendTo  []string
}
