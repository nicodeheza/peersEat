package modules

type Application struct {
	Peer *PeerModule
}

func InitApp() *Application {

	peerModule := newPeerModule()

	return &Application{peerModule}
}