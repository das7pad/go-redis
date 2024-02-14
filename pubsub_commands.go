package redis

import "context"

type PubSubCmdable interface {
	Publish(ctx context.Context, channel string, message interface{}) *IntCmd
	SPublish(ctx context.Context, channel string, message interface{}) *IntCmd
	PubSubChannels(ctx context.Context, pattern string) *StringSliceCmd
	PubSubNumSub(ctx context.Context, channels ...string) *MapStringIntCmd
	PubSubNumPat(ctx context.Context) *IntCmd
	PubSubShardChannels(ctx context.Context, pattern string) *StringSliceCmd
	PubSubShardNumSub(ctx context.Context, channels ...string) *MapStringIntCmd
}

// Publish posts the message to the channel.
func (c cmdable) Publish(ctx context.Context, channel string, message interface{}) *IntCmd {
	cmd := NewIntCmd2(ctx, "publish", channel, []interface{}{message})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) SPublish(ctx context.Context, channel string, message interface{}) *IntCmd {
	cmd := NewIntCmd2(ctx, "spublish", channel, []interface{}{message})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) PubSubChannels(ctx context.Context, pattern string) *StringSliceCmd {
	var argsS []string
	if pattern != "*" {
		argsS = append(argsS, pattern)
	}
	cmd := NewStringSliceCmd2S(ctx, "pubsub", "channels", argsS)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) PubSubNumSub(ctx context.Context, channels ...string) *MapStringIntCmd {
	cmd := NewMapStringIntCmd2S(ctx, "pubsub", "numsub", channels)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) PubSubShardChannels(ctx context.Context, pattern string) *StringSliceCmd {
	var argsS []string
	if pattern != "*" {
		argsS = append(argsS, pattern)
	}
	cmd := NewStringSliceCmd2S(ctx, "pubsub", "shardchannels", argsS)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) PubSubShardNumSub(ctx context.Context, channels ...string) *MapStringIntCmd {
	cmd := NewMapStringIntCmd2S(ctx, "pubsub", "shardnumsub", channels)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) PubSubNumPat(ctx context.Context) *IntCmd {
	cmd := NewIntCmd2(ctx, "pubsub", "numpat", nil)
	_ = c(ctx, cmd)
	return cmd
}
