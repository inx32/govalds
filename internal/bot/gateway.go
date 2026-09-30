package bot

import (
	"sync"

	"github.com/disgoorg/disgo/bot"
	"github.com/rs/zerolog"
)

type gatewayWrapper struct {
	next bot.EventListener
	log  zerolog.Logger

	wg     sync.WaitGroup
	closed bool
	mu     sync.RWMutex
}

func (w *gatewayWrapper) OnEvent(event bot.Event) {
	if w.next == nil {
		panic("handle gateway event: next handler is nil")
	}

	w.mu.RLock()
	if w.closed {
		w.mu.RUnlock()
		return
	}
	w.wg.Add(1)
	w.mu.RUnlock()

	go func() {
		defer w.wg.Done()
		w.next.OnEvent(event)
	}()
}

func (w *gatewayWrapper) Close() {
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()
		return
	}
	w.closed = true
	w.mu.Unlock()

	w.log.Debug().Msg("Waiting for running gateway event handlers to finish...")
	w.wg.Wait()
}
