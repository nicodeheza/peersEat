package mocks

import (
	"errors"
	"math"

	"github.com/nicodeheza/peersEat/constants"
	"github.com/nicodeheza/peersEat/models"
)

type GeoService struct{}

func NewGeo() *GeoService {
	return &GeoService{}
}

func (g *GeoService) GetCoordDistance(coord1 models.GeoCoords, coord2 models.GeoCoords) float64 {
	return coord1.Long
}

func (g *GeoService) IsSameCoord(coord1 models.GeoCoords, coord2 models.GeoCoords) bool {
	if coord1.Lat == coord2.Lat &&
		coord1.Long == coord2.Long {
		return true
	}

	return false
}

// update
func (g *GeoService) IsInInfluenceArea(peerCenter, geoPoint models.GeoCoords) bool {
	if g.IsSameCoord(peerCenter, geoPoint) {
		return true
	}

	peerDis := g.GetCoordDistance(peerCenter, geoPoint)
	peerDis = math.Abs(peerDis)

	if peerDis <= constants.INFLUENCE_RADIUS {
		return true
	}
	return false
}

func (g *GeoService) AreInfluenceAreasOverlaying(selfPeer models.Peer, peer models.Peer) bool {
	if peer.DeliveryRadius == 2 {
		return true
	}
	return false
}

func (g *GeoService) IsInDeliveryArea(selfPeer models.Peer, peer models.Peer) bool {
	if peer.DeliveryRadius == 1 {
		return true
	}
	return false
}

func (g *GeoService) GetAddressCoords(address, city, country string) (models.GeoCoords, error) {
	if address == "error" {
		return models.GeoCoords{}, errors.New("test error")
	}

	return models.GeoCoords{Long: 1.1, Lat: 2.2}, nil
}
