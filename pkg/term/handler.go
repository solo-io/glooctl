package term

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var terminationSignals = []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}

type Handler struct {
	notify []func()
	final  func(os.Signal)
	once   sync.Once
}

func Chain(handler *Handler, notify ...func()) *Handler {
	if handler == nil {
		return New(nil, notify...)
	}
	return New(handler.Signal, append(notify, handler.Close)...)
}

func New(final func(os.Signal), notify ...func()) *Handler {
	return &Handler{
		final:  final,
		notify: notify,
	}
}

func (h *Handler) Close() {
	h.once.Do(func() {
		for _, fn := range h.notify {
			fn()
		}
	})
}

func (h *Handler) Signal(s os.Signal) {
	h.once.Do(func() {
		for _, fn := range h.notify {
			fn()
		}
		if h.final == nil {
			os.Exit(1)
		}
		h.final(s)
	})
}

func (h *Handler) Run(fn func() error) error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, terminationSignals...)
	defer func() {
		signal.Stop(ch)
		close(ch)
	}()
	defer h.Close()
	return fn()
}
