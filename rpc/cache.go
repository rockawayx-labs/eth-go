package rpc

import "context"

type Cache interface {
	Set(ctx context.Context, key string, response []byte)
	Get(ctx context.Context, key string) (data []byte, found bool)
}
