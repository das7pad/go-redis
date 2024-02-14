package redis

import "context"

type HyperLogLogCmdable interface {
	PFAdd(ctx context.Context, key string, els ...interface{}) *IntCmd
	PFCount(ctx context.Context, keys ...string) *IntCmd
	PFMerge(ctx context.Context, dest string, keys ...string) *StatusCmd
}

func (c cmdable) PFAdd(ctx context.Context, key string, els ...interface{}) *IntCmd {
	cmd := NewIntCmd2Any(ctx, "pfadd", key, els)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) PFCount(ctx context.Context, keys ...string) *IntCmd {
	cmd := NewIntCmd2S(ctx, "pfcount", "", keys)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) PFMerge(ctx context.Context, dest string, keys ...string) *StatusCmd {
	cmd := NewStatusCmd2S(ctx, "pfmerge", dest, keys)
	_ = c(ctx, cmd)
	return cmd
}
