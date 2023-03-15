package mocks

import (
	"github.com/nicodeheza/peersEat/models"
)

type GeoService struct{}

func NewGeo() *GeoService{
	return &GeoService{}
}

func (g GeoService) GetCoorDistance(coor1 models.Center, coor2 models.Center) float64{
	return coor1.Long
}

func (g GeoService) IsSameCoor( coor1 models.Center, coor2 models.Center) bool {
	if coor1.Lat == coor2.Lat &&
	 coor1.Long == coor2.Long{
		return true
	}

	return false
}

func (g GeoService) IsInInfluenceArea(selfPeer models.Peer, peer models.Peer) bool{
	if peer.InfluenceRadius == 1{
		return true
	}
	return false
}

func (g GeoService) IsInDeliveryArea(selfPeer models.Peer, peer models.Peer) bool{
	if peer.DeliveryRadius == 1{
		return true
	}
	return false
}

