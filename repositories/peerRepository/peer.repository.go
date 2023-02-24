package peer_repository

import (
	"context"

	"github.com/nicodeheza/peersEat/models"
)



func Insert(peer models.Peer) error {
_,err := models.GetPeerColl().InsertOne(context.Background(), peer)
if err != nil{ return err}
return nil
}