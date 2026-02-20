package http

import "context"

type ICache interface {
	IsRequestLimiterExceeded(ctx context.Context, URI string) bool
}
