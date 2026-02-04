package pretty

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"strconv"
	"sync"
)

// Handler — slog.Handler с человекочитаемым выводом для локальной разработки.
type Handler struct {
	mu  sync.Mutex
	w   io.Writer
	lvl slog.Level
}

// NewHandler создаёт новый pretty handler.
func NewHandler(w io.Writer, level slog.Level) *Handler {
	return &Handler{w: w, lvl: level}
}

// Enabled reports whether the handler handles records at the given level.
func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.lvl
}

// Handle formats the Record and writes to h.w.
func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	b := make([]byte, 0, 256)
	b = append(b, r.Time.Format("15:04:05")...)
	b = append(b, ' ')
	b = append(b, levelShort(r.Level)...)
	b = append(b, ' ')
	b = append(b, r.Message...)

	r.Attrs(func(a slog.Attr) bool {
		b = append(b, ' ')
		b = append(b, a.Key...)
		b = append(b, '=')
		b = appendValue(b, a.Value)
		return true
	})
	b = append(b, '\n')

	_, err := h.w.Write(b)
	return err
}

func levelShort(l slog.Level) string {
	switch l {
	case slog.LevelDebug:
		return "DBG"
	case slog.LevelInfo:
		return "INF"
	case slog.LevelWarn:
		return "WRN"
	case slog.LevelError:
		return "ERR"
	}
	return "???"
}

func appendValue(b []byte, v slog.Value) []byte {
	switch v.Kind() {
	case slog.KindString:
		return append(b, v.String()...)
	case slog.KindInt64:
		return append(b, strconv.FormatInt(v.Int64(), 10)...)
	case slog.KindUint64:
		return append(b, strconv.FormatUint(v.Uint64(), 10)...)
	case slog.KindBool:
		if v.Bool() {
			return append(b, "true"...)
		}
		return append(b, "false"...)
	case slog.KindDuration:
		return append(b, v.Duration().String()...)
	case slog.KindAny:
		j, _ := json.Marshal(v.Any())
		return append(b, string(j)...)
	default:
		return append(b, v.String()...)
	}
}

// WithAttrs returns a new handler with the given attributes (Handler ignores them at the root).
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns a new handler for the given group (Handler ignores groups).
func (h *Handler) WithGroup(name string) slog.Handler {
	return h
}
