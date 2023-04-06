package geo

import (
	"math"
	"reflect"
	"testing"

	"github.com/nicodeheza/peersEat/models"
)

type GetCoorDistanceTests struct {
	Coor1  models.GeoCoords
	Coor2  models.GeoCoords
	result float64
}

func TestGetCoorDistance(t *testing.T) {
	g := NewGeo()
	tests := []GetCoorDistanceTests{
		{
			models.GeoCoords{Long: -36.0233352, Lat: -61.2301026},
			models.GeoCoords{Long: -36.1414237, Lat: -59.7854004},
			160.765761,
		},
		{
			models.GeoCoords{Long: 54.4828045, Lat: -312.1237794},
			models.GeoCoords{Long: 54.4904641, Lat: -311.7656207},
			39.827583,
		},
		{
			models.GeoCoords{Long: 54.8380309, Lat: -476.8637697},
			models.GeoCoords{Long: 54.8194413, Lat: -476.8588257},
			1.083700,
		},
	}

	for _, test := range tests {
		result := g.GetCoordDistance(test.Coor1, test.Coor2)
		dif := math.Abs(result - test.result)
		if dif > 0.001 {
			t.Errorf("expect %f but gets %f for %f, %f - %f, %f",
				test.result, result, test.Coor1.Lat, test.Coor1.Long,
				test.Coor2.Lat, test.Coor2.Long)
		}
	}
}

type IsSameCoorTest struct {
	Coor1  models.GeoCoords
	Coor2  models.GeoCoords
	result bool
}

func TestIsSameCoor(t *testing.T) {
	g := NewGeo()
	tests := []IsSameCoorTest{
		{
			models.GeoCoords{Long: -36.0233352, Lat: -61.2301026},
			models.GeoCoords{Long: -36.0233352, Lat: -59.7854004},
			false,
		},
		{
			models.GeoCoords{Long: 54.4828045, Lat: -312.1237794},
			models.GeoCoords{Long: 54.4828045, Lat: -312.1237794},
			true,
		},
	}
	for _, test := range tests {

		result := g.IsSameCoord(test.Coor1, test.Coor2)

		if result != test.result {
			t.Errorf("expect %t, but gets %t on (%v, %v)",
				test.result, result, test.Coor1, test.Coor2)
		}
	}

}

type IsInInfluenceAreaTest struct {
	SelfPeer models.Peer
	Peer     models.Peer
	result   bool
}

func TestIsInInfluenceArea(t *testing.T) {
	g := NewGeo()

	basePeer := models.Peer{
		Url:     "http://test.com",
		Center:  models.GeoCoords{Long: -36.0233352, Lat: -61.2301026},
		City:    "test city",
		Country: "test country",
	}
	peer1 := basePeer
	peer1.Center = models.GeoCoords{Long: -32.0233352, Lat: -61.2301026}

	peer2 := basePeer
	peer2.DeliveryRadius = 2
	peer2.Center = models.GeoCoords{Long: -34.605447, Lat: -58.383594}

	peer3 := peer2
	peer3.Center = models.GeoCoords{Long: -34.605341, Lat: -58.386598}

	peer4 := peer2
	peer4.Center = models.GeoCoords{Long: -34.606716, Lat: -58.470303}
	tests := []IsInInfluenceAreaTest{
		{
			basePeer,
			basePeer,
			true,
		},
		{
			basePeer,
			peer1,
			false,
		},
		{
			peer2,
			peer3,
			true,
		},
		{
			peer2,
			peer4,
			false,
		},
	}

	for _, test := range tests {
		result := g.AreInfluenceAreasOverlaying(test.SelfPeer, test.Peer)

		if result != test.result {
			t.Errorf("expect %t but gets %t for (%v, %v)", test.result, result, test.SelfPeer, test.Peer)
		}
	}
}

func TestIsInDeliveryArea(t *testing.T) {
	g := NewGeo()

	basePeer := models.Peer{
		Url:     "http://test.com",
		Center:  models.GeoCoords{Long: -36.0233352, Lat: -61.2301026},
		City:    "test city",
		Country: "test country",
	}
	peer1 := basePeer
	peer1.Center = models.GeoCoords{Long: -32.0233352, Lat: -61.2301026}

	peer2 := basePeer
	peer2.DeliveryRadius = 2
	peer2.Center = models.GeoCoords{Long: -34.605447, Lat: -58.383594}

	peer3 := peer2
	peer3.Center = models.GeoCoords{Long: -34.605341, Lat: -58.386598}

	peer4 := peer2
	peer4.Center = models.GeoCoords{Long: -34.606716, Lat: -58.470303}
	tests := []IsInInfluenceAreaTest{
		{
			basePeer,
			basePeer,
			true,
		},
		{
			basePeer,
			peer1,
			false,
		},
		{
			peer2,
			peer3,
			true,
		},
		{
			peer2,
			peer4,
			false,
		},
	}

	for _, test := range tests {
		result := g.IsInDeliveryArea(test.SelfPeer, test.Peer)

		if result != test.result {
			t.Errorf("expect %t but gets %t for (%v, %v)", test.result, result, test.SelfPeer, test.Peer)
		}
	}
}

func TestGetAddressCords(t *testing.T) {
	g := NewGeo()

	res, err := g.GetAddressCoords("cerrito 800", "Buenos Aires", "Argentina")
	if err != nil {
		t.Error(err.Error())
	}

	expected := models.GeoCoords{Long: -58.3826948, Lat: -34.5992499}

	if !reflect.DeepEqual(expected, res) {
		t.Errorf("Expected: %v but go: %v", expected, res)
	}
}
