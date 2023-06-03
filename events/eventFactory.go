package events

import (
	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/types"
)

const ADD_NEW_PEER = "addPeer"
const DELIVERY_AREA_UPDATED = "deliveryAreaUpdated"

func NewAddPeerEvent(peer models.Peer, sendTo []string) types.Event {
	return types.Event{
		Name:    ADD_NEW_PEER,
		Payload: peer,
		SendTo:  sendTo,
	}
}

func NewUpdateDeliveryAreaEvent(peer models.Peer, sendTo []string) types.Event {
	return types.Event{
		Name:    DELIVERY_AREA_UPDATED,
		Payload: peer,
		SendTo:  sendTo,
	}
}
