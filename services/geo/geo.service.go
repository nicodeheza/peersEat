package geo

import (
	"math"

	"github.com/nicodeheza/peersEat/models"
)

// do it injectable

// https://gist.github.com/hotdang-ca/6c1ee75c48e515aec5bc6db6e3265e49
func GetCoorDistance(coor1 models.Center, coor2 models.Center) float64{
	const R= 6371e3
	radlat1 := float64(math.Pi * coor1.Lat / 180)
	radlat2 := float64(math.Pi * coor2.Lat / 180)

	theta := float64(coor1.Long - coor2.Long)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radlat1) * math.Sin(radlat2) + math.Cos(radlat1) * math.Cos(radlat2) * math.Cos(radtheta);
	if dist > 1{
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist *180 / math.Pi
	dist= dist * 60 * 1.1515

	return dist * 1.609344
}

func IsSameCoor( coor1 models.Center, coor2 models.Center) bool {
	if coor1.Lat == coor2.Lat &&
	 coor1.Long == coor2.Long{
		return true
	}

	return false
}

func IsInInfluenceArea(selfPeer models.Peer, peer models.Peer) bool{
	if IsSameCoor(selfPeer.Center, peer.Center){return true}

	if selfPeer.InfluenceRadius == 0 {return false}

	peerDis := GetCoorDistance(peer.Center, selfPeer.Center)
	peerDis = math.Abs(peerDis) 

	influenceSum := selfPeer.InfluenceRadius + peer.InfluenceRadius

	if peerDis <= influenceSum{ return true}
	return false
}

func IsInDeliveryArea(selfPeer models.Peer, peer models.Peer) bool{
	if IsSameCoor(selfPeer.Center, peer.Center){return true}

	if selfPeer.DeliveryRadius == 0 {return false}

	peerDis := GetCoorDistance(peer.Center, selfPeer.Center)
	peerDis = math.Abs(peerDis) 

	deliverySum := selfPeer.DeliveryRadius + peer.DeliveryRadius

	if peerDis <= deliverySum{ return true}
	return false
}

