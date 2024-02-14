package redis

import (
	"context"
	"strconv"
	"strings"
	"time"
)

type ListCmdable interface {
	BLPop(ctx context.Context, timeout time.Duration, keys ...string) *StringSliceCmd
	BLMPop(ctx context.Context, timeout time.Duration, direction string, count int64, keys ...string) *KeyValuesCmd
	BRPop(ctx context.Context, timeout time.Duration, keys ...string) *StringSliceCmd
	BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *StringCmd
	LIndex(ctx context.Context, key string, index int64) *StringCmd
	LInsert(ctx context.Context, key, op string, pivot, value interface{}) *IntCmd
	LInsertBefore(ctx context.Context, key string, pivot, value interface{}) *IntCmd
	LInsertAfter(ctx context.Context, key string, pivot, value interface{}) *IntCmd
	LLen(ctx context.Context, key string) *IntCmd
	LMPop(ctx context.Context, direction string, count int64, keys ...string) *KeyValuesCmd
	LPop(ctx context.Context, key string) *StringCmd
	LPopCount(ctx context.Context, key string, count int) *StringSliceCmd
	LPos(ctx context.Context, key string, value string, args LPosArgs) *IntCmd
	LPosCount(ctx context.Context, key string, value string, count int64, args LPosArgs) *IntSliceCmd
	LPush(ctx context.Context, key string, values ...interface{}) *IntCmd
	LPushX(ctx context.Context, key string, values ...interface{}) *IntCmd
	LRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd
	LRem(ctx context.Context, key string, count int64, value interface{}) *IntCmd
	LSet(ctx context.Context, key string, index int64, value interface{}) *StatusCmd
	LTrim(ctx context.Context, key string, start, stop int64) *StatusCmd
	RPop(ctx context.Context, key string) *StringCmd
	RPopCount(ctx context.Context, key string, count int) *StringSliceCmd
	RPopLPush(ctx context.Context, source, destination string) *StringCmd
	RPush(ctx context.Context, key string, values ...interface{}) *IntCmd
	RPushX(ctx context.Context, key string, values ...interface{}) *IntCmd
	LMove(ctx context.Context, source, destination, srcpos, destpos string) *StringCmd
	BLMove(ctx context.Context, source, destination, srcpos, destpos string, timeout time.Duration) *StringCmd
}

func (c cmdable) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *StringSliceCmd {
	argsS := append(keys, strconv.FormatInt(formatSec(ctx, timeout), 10))
	cmd := NewStringSliceCmd2S(ctx, "blpop", "", argsS)
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) BLMPop(ctx context.Context, timeout time.Duration, direction string, count int64, keys ...string) *KeyValuesCmd {
	args := make([]interface{}, 2+len(keys), 5+len(keys))
	args[0] = formatSec(ctx, timeout)
	args[1] = len(keys)
	for i, key := range keys {
		args[i+2] = key
	}
	args = append(args, strings.ToLower(direction), "count", count)
	cmd := NewKeyValuesCmd2(ctx, "blmpop", "", args)
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *StringSliceCmd {
	argsS := append(keys, strconv.FormatInt(formatSec(ctx, timeout), 10))
	cmd := NewStringSliceCmd2S(ctx, "brpop", "", argsS)
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *StringCmd {
	cmd := NewStringCmd3(
		ctx,
		"brpoplpush",
		source,
		destination,
		[]interface{}{formatSec(ctx, timeout)},
	)
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) LIndex(ctx context.Context, key string, index int64) *StringCmd {
	cmd := NewStringCmd2(ctx, "lindex", key, []interface{}{index})
	_ = c(ctx, cmd)
	return cmd
}

// LMPop Pops one or more elements from the first non-empty list key from the list of provided key names.
// direction: left or right, count: > 0
// example: client.LMPop(ctx, "left", 3, "key1", "key2")
func (c cmdable) LMPop(ctx context.Context, direction string, count int64, keys ...string) *KeyValuesCmd {
	args := make([]interface{}, 1+len(keys), 4+len(keys))
	args[0] = len(keys)
	for i, key := range keys {
		args[i+1] = key
	}
	args = append(args, strings.ToLower(direction), "count", count)
	cmd := NewKeyValuesCmd2(ctx, "lmpop", "", args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *IntCmd {
	cmd := NewIntCmd3(ctx, "linsert", key, op, []interface{}{pivot, value})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) LInsertBefore(ctx context.Context, key string, pivot, value interface{}) *IntCmd {
	cmd := NewIntCmd3(ctx, "linsert", key, "before", []interface{}{pivot, value})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) LInsertAfter(ctx context.Context, key string, pivot, value interface{}) *IntCmd {
	cmd := NewIntCmd3(ctx, "linsert", key, "after", []interface{}{pivot, value})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) LLen(ctx context.Context, key string) *IntCmd {
	cmd := NewIntCmd2(ctx, "llen", key, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) LPop(ctx context.Context, key string) *StringCmd {
	cmd := NewStringCmd2(ctx, "lpop", key, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) LPopCount(ctx context.Context, key string, count int) *StringSliceCmd {
	cmd := NewStringSliceCmd2(ctx, "lpop", key, []interface{}{count})
	_ = c(ctx, cmd)
	return cmd
}

type LPosArgs struct {
	Rank, MaxLen int64
}

func (c cmdable) LPos(ctx context.Context, key string, value string, a LPosArgs) *IntCmd {
	var args []interface{}
	if a.Rank != 0 {
		args = append(args, "rank", a.Rank)
	}
	if a.MaxLen != 0 {
		args = append(args, "maxlen", a.MaxLen)
	}

	cmd := NewIntCmd3(ctx, "lpos", key, value, args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) LPosCount(ctx context.Context, key string, value string, count int64, a LPosArgs) *IntSliceCmd {
	args := []interface{}{value, "count", count}
	if a.Rank != 0 {
		args = append(args, "rank", a.Rank)
	}
	if a.MaxLen != 0 {
		args = append(args, "maxlen", a.MaxLen)
	}
	cmd := NewIntSliceCmd2(ctx, "lpos", key, args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) LPush(ctx context.Context, key string, values ...interface{}) *IntCmd {
	cmd := NewIntCmd2Any(ctx, "lpush", key, values)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) LPushX(ctx context.Context, key string, values ...interface{}) *IntCmd {
	cmd := NewIntCmd2Any(ctx, "lpushx", key, values)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) LRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd {
	cmd := NewStringSliceCmd2(
		ctx,
		"lrange",
		key,
		[]interface{}{start, stop},
	)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) LRem(ctx context.Context, key string, count int64, value interface{}) *IntCmd {
	cmd := NewIntCmd2(ctx, "lrem", key, []interface{}{count, value})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) LSet(ctx context.Context, key string, index int64, value interface{}) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "lset", key, []interface{}{index, value})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) LTrim(ctx context.Context, key string, start, stop int64) *StatusCmd {
	cmd := NewStatusCmd2(
		ctx,
		"ltrim",
		key,
		[]interface{}{start, stop},
	)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) RPop(ctx context.Context, key string) *StringCmd {
	cmd := NewStringCmd2(ctx, "rpop", key, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) RPopCount(ctx context.Context, key string, count int) *StringSliceCmd {
	cmd := NewStringSliceCmd2(ctx, "rpop", key, []interface{}{count})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) RPopLPush(ctx context.Context, source, destination string) *StringCmd {
	cmd := NewStringCmd3(ctx, "rpoplpush", source, destination, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) RPush(ctx context.Context, key string, values ...interface{}) *IntCmd {
	cmd := NewIntCmd2Any(ctx, "rpush", key, values)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) RPushX(ctx context.Context, key string, values ...interface{}) *IntCmd {
	cmd := NewIntCmd2Any(ctx, "rpushx", key, values)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) LMove(ctx context.Context, source, destination, srcpos, destpos string) *StringCmd {
	cmd := NewStringCmd3S(ctx, "lmove", source, destination, []string{srcpos, destpos})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) BLMove(ctx context.Context, source, destination, srcpos, destpos string, timeout time.Duration) *StringCmd {
	cmd := NewStringCmd3(ctx, "blmove", source, destination, []interface{}{srcpos, destpos, formatSec(ctx, timeout)})
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
}
