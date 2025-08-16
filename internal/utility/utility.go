package utility

import (
	"context"
	"time"
)

type ContextType int

const (
	DB ContextType = iota
	HTTP
)

const (
	ctxDBDuration      int = 3
	ctxHTTPDuration    int = 5
	ctxDefaultDuration int = 10
)

func GenerateContextWithTimeout(ctxType ContextType) (context.Context, context.CancelFunc) {
	var duration time.Duration

	switch ctxType {
	case DB:
		duration = time.Duration(ctxDBDuration)
	case HTTP:
		duration = time.Duration(ctxHTTPDuration)
	default:
		duration = time.Duration(ctxDefaultDuration)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*duration)

	return ctx, cancel
}
