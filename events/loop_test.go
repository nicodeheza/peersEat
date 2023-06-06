package events

import (
	"container/list"
	"testing"

	"github.com/nicodeheza/peersEat/models"
)

func TestLoop(t *testing.T) {
	queue := list.New()
	loop := EventLoop{queue: queue}

	event := NewAddPeerEvent(models.Peer{Url: "http://tests.com"}, []string{})
	loop.Enqueue(event)
	if loop.queueLen() != 1 {
		t.Error("expect queue have len 1")
	}

	result := loop.Dequeue()
	if result.Name != event.Name {
		t.Errorf("expect %s to be equal to %s", result.Name, event.Name)
	}
	if result.Payload.(models.Peer).Url != event.Payload.(models.Peer).Url {
		t.Errorf("expect %s to be equal to %s", result.Payload.(models.Peer).Url, event.Payload.(models.Peer).Url)
	}
	if loop.queueLen() != 0 {
		t.Errorf("expect queue have 0 len, but have %d", loop.queueLen())
	}
}
