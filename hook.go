package rotatelogs

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
)

func NewHook(path string, options ...Option) (log.Hook, error) {
	var err error
	h := &hook{}

	for _, option := range options {
		if option.Name() == optkeyFormatter {
			if formatter, ok := option.Value().(log.Formatter); ok {
				h.formatter = formatter
				break
			}
		}
	}

	if h.out, err = New(path, options...); err != nil {
		return nil, err
	}

	if h.formatter == nil {
		h.formatter = NewTextFormatter(true)
	}

	return h, nil
}

type hook struct {
	formatter log.Formatter
	out       io.Writer
}

var (
	errHookInitial = errors.New("rotate logs hook is not initialized")
)

// Fire writes the log file to defined path or using the defined writer.
func (h *hook) Fire(entry *log.Entry) error {
	if h.formatter == nil || h.out == nil {
		return errHookInitial
	}

	msg, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = h.out.Write(msg)
	return err
}

// Levels returns configured log levels.
func (h *hook) Levels() []log.Level {
	return log.AllLevels
}
