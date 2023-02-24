package peer_service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nicodeheza/peersEat/config"
	"github.com/nicodeheza/peersEat/models"
	peer_repository "github.com/nicodeheza/peersEat/repositories/peerRepository"
)

func InitPeer() {
	centerSrt:= config.GetEnv("CENTER")
	centerSlice := strings.Split(centerSrt, ",")
	long, _ := strconv.ParseFloat(centerSlice[0], 64)
	lat, _ := strconv.ParseFloat(centerSlice[1], 64)

	selfPeer := models.Peer{
		Url: config.GetEnv("HOST"),
		Center: models.Center{Long: long , Lat: lat},
		City: config.GetEnv("CITY"),
		Country: config.GetEnv("COUNTRY"),
	}

	fmt.Println(selfPeer)

	peer_repository.Insert(selfPeer)
}