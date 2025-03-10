package mylog

import (
	"context"
	"fmt"
	"github.com/lmittmann/tint"
	"io"
	"log/slog"
	"sync"
	"time"
)

type (
	customConsoleHandler struct {
		slog.Handler

		w  io.Writer
		mu sync.Mutex
	}
)

func (h *customConsoleHandler) Handle(ctx context.Context, record slog.Record) error {
	newRecord := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)

	var err error
	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key == errKey {
			if e, ok := attr.Value.Any().(error); !ok {
				return true
			} else {
				err = e
			}
			newRecord.AddAttrs(tint.Err(err))
		} else {
			newRecord.AddAttrs(attr)
		}

		return true
	})

	if err := h.Handler.Handle(ctx, newRecord); err != nil {
		return err
	}

	if err != nil {
		h.mu.Lock()
		defer h.mu.Unlock()

		if _, err := fmt.Fprintf(h.w, "%+v\n", err); err != nil {
			return err
		}
	}

	return nil
}

func newHandler(logLevel slog.Level, w io.Writer) slog.Handler {
	handler := tint.NewHandler(w, &tint.Options{
		AddSource:  true,
		NoColor:    false,
		TimeFormat: time.RFC3339,
		Level:      logLevel,
	})

	return &customConsoleHandler{
		Handler: handler,
		w:       w,
	}
}

const errKey = "err"

func Err(err error) slog.Attr {
	return slog.Any(errKey, err)
}
