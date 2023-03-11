package types

import "github.com/nicodeheza/peersEat/models"

type PeerPresentationBody struct {
	NewPeer models.Peer
	SendTo  []string
}

type SendAllPeerQuery struct{
	Excludes []string `query:"excludes"`
}