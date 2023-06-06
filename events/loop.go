package events

import (
	"container/list"
	"sync"
	"time"

	"github.com/nicodeheza/peersEat/types"
)

type EventLoop struct {
	mutex    sync.Mutex
	queue    *list.List
	handlers HandlersI
}

type EventLoopI interface {
	Enqueue(event types.Event)
	Dequeue() types.Event
	PropagateEvent(event types.Event)
}

var eventLoopInstance *EventLoop
var once sync.Once

func InitEventLoop(handlers HandlersI) *EventLoop {
	once.Do(func() {
		queue := list.New()
		eventLoop := EventLoop{queue: queue, handlers: handlers}
		go eventLoop.Loop()
		eventLoopInstance = &eventLoop
	})
	return eventLoopInstance
}

func (e *EventLoop) Enqueue(event types.Event) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.queue.PushBack(event)
}

func (e *EventLoop) Dequeue() types.Event {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	element := e.queue.Front()
	e.queue.Remove(element)

	event, ok := element.Value.(types.Event)
	if !ok {
		return types.Event{}
	}
	return event
}

func (e *EventLoop) queueLen() int {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.queue.Len()
}

func (e *EventLoop) PropagateEvent(event types.Event) {
	e.handlers.PropagateEvent(event)
}

func (e *EventLoop) Loop() {
	for {
		time.Sleep(1 * time.Second)

		if e.queueLen() == 0 {
			continue
		}

		event := e.Dequeue()

		switch event.Name {
		case ADD_NEW_PEER:
			e.handlers.HandleAddPeer(event)
		case DELIVERY_AREA_UPDATED:
			e.handlers.PeerUpdatedDeliveryArea(event)
		default:
			continue
		}
	}
}
