package peer_service

import (
	"os"
	"strconv"
	"strings"

	"github.com/nicodeheza/peersEat/models"
	peer_repository "github.com/nicodeheza/peersEat/repositories/peerRepository"
)

func InitPeer() {
	centerSrt:= os.Getenv("CENTER")
	centerSlice := strings.Split(centerSrt, ",")
	long, _ := strconv.ParseFloat(centerSlice[0], 64)
	lat, _ := strconv.ParseFloat(centerSlice[1], 64)

	selfPeer := models.Peer{
		Url: os.Getenv("HOST"),
		Center: models.Center{Long: long , Lat: lat},
		City: os.Getenv("CITY"),
		Country: os.Getenv("COUNTRY"),
	}

	peer_repository.Insert(selfPeer)

	// Todo: create new peers endpoint and add here call 
}