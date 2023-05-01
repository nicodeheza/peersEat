package modules

import (
	"github.com/nicodeheza/peersEat/events"
	"github.com/nicodeheza/peersEat/repositories"
	"github.com/nicodeheza/peersEat/services/geo"
	"github.com/nicodeheza/peersEat/services/validations"
)

type EventModule struct {
	loop     events.EventLoopI
	handlers events.HandlersI
}

func InitEventModules(
	peerRepo repositories.PeerRepositoryI,
	validation validations.ValidateI,
	geo geo.GeoServiceI,
) *EventModule {
	handlers := events.NewEventHandlers(peerRepo, validation, geo)
	loop := events.InitEventLoop(handlers)
	return &EventModule{
		loop, handlers,
	}
}
