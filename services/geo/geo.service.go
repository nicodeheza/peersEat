package geo

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/nicodeheza/peersEat/constants"
	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/types"
)

type GeoService struct{}
type GeoServiceI interface {
	GetCoordDistance(coord1 models.GeoCoords, coord2 models.GeoCoords) float64
	IsSameCoord(coord1 models.GeoCoords, coord2 models.GeoCoords) bool
	IsInInfluenceArea(peerCenter, geoPoint models.GeoCoords) bool
	AreInfluenceAreasOverlaying(selfPeer models.Peer, peer models.Peer) bool
	IsInDeliveryArea(selfPeer models.Peer, peer models.Peer) bool
	GetAddressCoords(address, city, country string) (models.GeoCoords, error)
}

func NewGeo() *GeoService {
	return &GeoService{}
}

// https://gist.github.com/hotdang-ca/6c1ee75c48e515aec5bc6db6e3265e49
func (g *GeoService) GetCoordDistance(coord1 models.GeoCoords, coord2 models.GeoCoords) float64 {
	const R = 6371e3
	radlat1 := float64(math.Pi * coord1.Lat / 180)
	radlat2 := float64(math.Pi * coord2.Lat / 180)

	theta := float64(coord1.Long - coord2.Long)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515

	return dist * 1.609344
}

func (g *GeoService) IsSameCoord(coord1 models.GeoCoords, coord2 models.GeoCoords) bool {
	if coord1.Lat == coord2.Lat &&
		coord1.Long == coord2.Long {
		return true
	}

	return false
}

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
	if g.IsSameCoord(selfPeer.Center, peer.Center) {
		return true
	}

	peerDis := g.GetCoordDistance(peer.Center, selfPeer.Center)
	peerDis = math.Abs(peerDis)

	if peerDis <= constants.INFLUENCE_RADIUS*2 {
		return true
	}
	return false
}

func (g *GeoService) IsInDeliveryArea(selfPeer models.Peer, peer models.Peer) bool {
	if g.IsSameCoord(selfPeer.Center, peer.Center) {
		return true
	}

	if selfPeer.DeliveryRadius == 0 {
		return false
	}

	peerDis := g.GetCoordDistance(peer.Center, selfPeer.Center)
	peerDis = math.Abs(peerDis)

	deliverySum := selfPeer.DeliveryRadius + peer.DeliveryRadius

	if peerDis <= deliverySum {
		return true
	}
	return false
}

func (g *GeoService) GetAddressCoords(address, city, country string) (models.GeoCoords, error) {
	url, err := url.Parse("https://nominatim.openstreetmap.org/search")
	if err != nil {
		return models.GeoCoords{}, err
	}

	query := url.Query()
	query.Add("q", fmt.Sprintf("%s,%s,%s", address, city, country))
	query.Add("format", "json")
	query.Add("limit", "1")
	url.RawQuery = query.Encode()

	resp, err := http.Get(url.String())
	if err != nil {
		return models.GeoCoords{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var getCordsResponse []types.GetCordsResponse
	err = decoder.Decode(&getCordsResponse)
	if err != nil {
		return models.GeoCoords{}, err
	}

	long, err := strconv.ParseFloat(getCordsResponse[0].Lon, 64)
	if err != nil {
		return models.GeoCoords{}, err
	}
	lat, err := strconv.ParseFloat(getCordsResponse[0].Lat, 64)
	if err != nil {
		return models.GeoCoords{}, err
	}

	return models.GeoCoords{Long: long, Lat: lat}, nil
}
