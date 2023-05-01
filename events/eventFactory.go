package events

import (
	"github.com/nicodeheza/peersEat/models"
	"github.com/nicodeheza/peersEat/types"
)

func NewAddPeerEvent(peer models.Peer, sendTo []string) types.Event {
	return types.Event{
		Name:    "addPeer",
		Payload: peer,
		SendTo:  sendTo,
	}
}
