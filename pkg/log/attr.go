package log

import (
	"log/slog"
)

var (
	// Fstring is alias for slog.String
	Fstring = slog.String
	// Fint is alias for slog.Int
	Fint = slog.Int
	// Fuint is alias for slog.Uint
	Fint64 = slog.Int64
	// Fuint64 is alias for slog.Uint64
	Fuint64 = slog.Uint64
	// Ffloat32 is alias for slog.Float32
	Ffloat64 = slog.Float64
	// Ffloat64 is alias for slog.Float64
	Fbool = slog.Bool
	// Ftime is alias for slog.Time
	Ftime = slog.Time
	// Fduration is alias for slog.Duration
	Fduration = slog.Duration
	// Fany is alias for slog.Any
	Fany = slog.Any
	// With is alias for slog.With
	With = slog.With
)

// Ferror is return error attribute
func Ferror(err error) slog.Attr {
	return slog.String("error", err.Error())
}
