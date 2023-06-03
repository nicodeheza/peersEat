package mocks

import "github.com/nicodeheza/peersEat/types"

type EventLoopMock struct{}

func NewEventLoopMock() *EventLoopMock {
	return &EventLoopMock{}
}

func (e *EventLoopMock) Enqueue(event types.Event) {

}

func (e *EventLoopMock) Dequeue() types.Event {
	return types.Event{}
}

func (e *EventLoopMock) PropagateEvent(event types.Event) {}
