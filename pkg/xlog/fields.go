package xlog

import "log/slog"

type Fields map[string]interface{}

func (f Fields) ToMap() map[string]interface{} {
	return map[string]interface{}(f)
}

func (f Fields) ToAttrs() []slog.Attr {
	attrs := make([]slog.Attr, 0, len(f))
	for key, value := range f {
		attrs = append(attrs, slog.Any(key, value))
	}
	return attrs
}
