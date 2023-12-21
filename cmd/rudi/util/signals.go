// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package util

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

type SignalHandler struct {
	queue  chan os.Signal
	cancel context.CancelFunc
	lock   sync.Mutex
}

func SetupSignalHandler() *SignalHandler {
	osSignals := make(chan os.Signal, 2)
	signal.Notify(osSignals, os.Interrupt)

	handler := &SignalHandler{
		queue: osSignals,
		lock:  sync.Mutex{},
	}

	go handler.ProcessQueue()

	return handler
}

func (h *SignalHandler) ProcessQueue() {
	for signal := range h.queue {
		h.handleSignal(signal)
	}
}

func (h *SignalHandler) SetCancelFn(cancel context.CancelFunc) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.cancel = cancel
}

func (h *SignalHandler) handleSignal(_ os.Signal) {
	h.lock.Lock()
	defer h.lock.Unlock()

	// If there is something to interrupt, do so.
	if h.cancel != nil {
		h.cancel()
		h.cancel = nil

		return
	}

	// If even stopping the script takes too long, a second interrupt
	// will lead to an immediate shutdown.
	// Note that in console mode, when sitting at the prompt at not running
	// a script, the term reader package handles any possible Ctrl-C key presses.
	os.Exit(1) //nolint:gocritic
}
