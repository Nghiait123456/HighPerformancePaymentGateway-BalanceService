package main

import (
	"sync"
	"sync/atomic"
)

const valueStop = 1

type emergencyStop struct {
	stop uint32
	mu   sync.Mutex
}

type emergencyStopInterface interface {
	init()
	IsStop() bool
	ThrowEmergencyStop()
}

func (e *emergencyStop) init() {
	e.stop = 0
}

func (e *emergencyStop) IsStop() bool {
	return atomic.CompareAndSwapUint32(&e.stop, valueStop, valueStop)
}

func (e *emergencyStop) ThrowEmergencyStop() {
	e.mu.Lock()
	e.stop = valueStop
	e.mu.Unlock()
}

func NewEmergencyStop() emergencyStopInterface {
	e := emergencyStop{}
	e.init()
	return &e
}
