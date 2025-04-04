package log

import (
	"io"

	"github.com/rs/zerolog"
	"github.com/xlgmokha/x/pkg/convert"
)

func New(writer io.Writer, fields Fields) *zerolog.Logger {
	return convert.ToPtr(
		zerolog.
			New(writer).
			With().
			Timestamp().
			Fields(fields.ToMap()).
			Logger(),
	)
}
