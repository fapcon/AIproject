// Code generated by protoc-gen-dispatcher. DO NOT EDIT.
// protoc version: v5.26.0
// source: internal_events.proto

package intev

import "studentgit.kata.academy/quant/torque/pkg/logster"

type AckableInternalEvent struct {
	InternalEvent *InternalEvent
	Ack           func()
}

type AckableSubscribeInstruments struct {
	SubscribeInstruments *SubscribeInstruments
	Ack                  func()
}

type AckableNewOrderRequested struct {
	NewOrderRequested *NewOrderRequested
	Ack               func()
}

type InternalEventDispatcher struct {
	logger               logster.Logger
	started              bool
	subscribeInstruments chan AckableSubscribeInstruments
	newOrderRequested    chan AckableNewOrderRequested
}

func NewInternalEventDispatcher(logger logster.Logger) *InternalEventDispatcher {
	return &InternalEventDispatcher{
		logger: logger,
	}
}

func (d *InternalEventDispatcher) Run(ch <-chan AckableInternalEvent) {
	d.started = true
	for internalEvent := range ch {
		switch event := internalEvent.InternalEvent.Event.(type) {
		case *InternalEvent_SubscribeInstruments:
			if d.subscribeInstruments != nil {
				d.subscribeInstruments <- AckableSubscribeInstruments{
					SubscribeInstruments: event.SubscribeInstruments,
					Ack:                  internalEvent.Ack,
				}
			} else {
				internalEvent.Ack()
			}
		case *InternalEvent_NewOrderRequested:
			if d.newOrderRequested != nil {
				d.newOrderRequested <- AckableNewOrderRequested{
					NewOrderRequested: event.NewOrderRequested,
					Ack:               internalEvent.Ack,
				}
			} else {
				internalEvent.Ack()
			}
		default:
			d.logger.Errorf("unknown type %T", event)
		}
	}

	if d.subscribeInstruments != nil {
		close(d.subscribeInstruments)
	}

	if d.newOrderRequested != nil {
		close(d.newOrderRequested)
	}
}

func (d *InternalEventDispatcher) SubscribeInstruments() <-chan AckableSubscribeInstruments {
	if d.subscribeInstruments != nil {
		panic("SubscribeInstruments already called")
	}
	if d.started {
		panic("SubscribeInstruments called after Run")
	}
	d.subscribeInstruments = make(chan AckableSubscribeInstruments)
	return d.subscribeInstruments
}

func (d *InternalEventDispatcher) NewOrderRequested() <-chan AckableNewOrderRequested {
	if d.newOrderRequested != nil {
		panic("NewOrderRequested already called")
	}
	if d.started {
		panic("NewOrderRequested called after Run")
	}
	d.newOrderRequested = make(chan AckableNewOrderRequested)
	return d.newOrderRequested
}
