package term

import (
	"io"
	"os"

	"github.com/docker/docker/pkg/term"
)

type SafeFunc func() error

type TTY struct {
	In     io.Reader
	Out    io.Writer
	Raw    bool
	TryDev bool
	Parent *Handler
}

func (t TTY) IsTerminalIn() bool {
	return IsTerminal(t.In)
}

func (t TTY) IsTerminalOut() bool {
	return IsTerminal(t.Out)
}

func IsTerminal(i interface{}) bool {
	_, terminal := term.GetFdInfo(i)
	return terminal
}

func (t TTY) Safe(fn SafeFunc) error {
	inFd, isTerminal := term.GetFdInfo(t.In)

	if !isTerminal && t.TryDev {
		if f, err := os.Open("/dev/tty"); err == nil {
			defer f.Close()
			inFd = f.Fd()
			isTerminal = term.IsTerminal(inFd)
		}
	}
	if !isTerminal {
		return fn()
	}

	var state *term.State
	var err error
	if t.Raw {
		state, err = term.MakeRaw(inFd)
	} else {
		state, err = term.SaveState(inFd)
	}
	if err != nil {
		return err
	}
	return Chain(t.Parent, func() {
		term.RestoreTerminal(inFd, state)
	}).Run(fn)
}
