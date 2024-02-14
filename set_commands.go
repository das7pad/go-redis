package redis

import "context"

type SetCmdable interface {
	SAdd(ctx context.Context, key string, members ...interface{}) *IntCmd
	SCard(ctx context.Context, key string) *IntCmd
	SDiff(ctx context.Context, keys ...string) *StringSliceCmd
	SDiffStore(ctx context.Context, destination string, keys ...string) *IntCmd
	SInter(ctx context.Context, keys ...string) *StringSliceCmd
	SInterCard(ctx context.Context, limit int64, keys ...string) *IntCmd
	SInterStore(ctx context.Context, destination string, keys ...string) *IntCmd
	SIsMember(ctx context.Context, key string, member interface{}) *BoolCmd
	SMIsMember(ctx context.Context, key string, members ...interface{}) *BoolSliceCmd
	SMembers(ctx context.Context, key string) *StringSliceCmd
	SMembersMap(ctx context.Context, key string) *StringStructMapCmd
	SMove(ctx context.Context, source, destination string, member interface{}) *BoolCmd
	SPop(ctx context.Context, key string) *StringCmd
	SPopN(ctx context.Context, key string, count int64) *StringSliceCmd
	SRandMember(ctx context.Context, key string) *StringCmd
	SRandMemberN(ctx context.Context, key string, count int64) *StringSliceCmd
	SRem(ctx context.Context, key string, members ...interface{}) *IntCmd
	SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd
	SUnion(ctx context.Context, keys ...string) *StringSliceCmd
	SUnionStore(ctx context.Context, destination string, keys ...string) *IntCmd
}

// ------------------------------------------------------------------------------

func (c cmdable) SAdd(ctx context.Context, key string, members ...interface{}) *IntCmd {
	cmd := NewIntCmd2Any(ctx, "sadd", key, members)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) SCard(ctx context.Context, key string) *IntCmd {
	cmd := NewIntCmd2(ctx, "scard", key, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) SDiff(ctx context.Context, keys ...string) *StringSliceCmd {
	cmd := NewStringSliceCmd2S(ctx, "sdiff", "", keys)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) SDiffStore(ctx context.Context, destination string, keys ...string) *IntCmd {
	cmd := NewIntCmd2S(ctx, "sdiffstore", destination, keys)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) SInter(ctx context.Context, keys ...string) *StringSliceCmd {
	cmd := NewStringSliceCmd2S(ctx, "sinter", "", keys)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) SInterCard(ctx context.Context, limit int64, keys ...string) *IntCmd {
	args := make([]interface{}, 1+len(keys), 1+len(keys)+2)
	args[0] = len(keys)
	for i, key := range keys {
		args[i+1] = key
	}
	args = append(args, "limit", limit)
	cmd := NewIntCmd2(ctx, "sintercard", "", args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) SInterStore(ctx context.Context, destination string, keys ...string) *IntCmd {
	cmd := NewIntCmd2S(ctx, "sinterstore", destination, keys)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) SIsMember(ctx context.Context, key string, member interface{}) *BoolCmd {
	cmd := NewBoolCmd2(ctx, "sismember", key, []interface{}{member})
	_ = c(ctx, cmd)
	return cmd
}

// SMIsMember Redis `SMISMEMBER key member [member ...]` command.
func (c cmdable) SMIsMember(ctx context.Context, key string, members ...interface{}) *BoolSliceCmd {
	cmd := NewBoolSliceCmd2Any(ctx, "smismember", key, members)
	_ = c(ctx, cmd)
	return cmd
}

// SMembers Redis `SMEMBERS key` command output as a slice.
func (c cmdable) SMembers(ctx context.Context, key string) *StringSliceCmd {
	cmd := NewStringSliceCmd2(ctx, "smembers", key, nil)
	_ = c(ctx, cmd)
	return cmd
}

// SMembersMap Redis `SMEMBERS key` command output as a map.
func (c cmdable) SMembersMap(ctx context.Context, key string) *StringStructMapCmd {
	cmd := NewStringStructMapCmd2(ctx, "smembers", key, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) SMove(ctx context.Context, source, destination string, member interface{}) *BoolCmd {
	cmd := NewBoolCmd3(ctx, "smove", source, destination, []interface{}{member})
	_ = c(ctx, cmd)
	return cmd
}

// SPop Redis `SPOP key` command.
func (c cmdable) SPop(ctx context.Context, key string) *StringCmd {
	cmd := NewStringCmd2(ctx, "spop", key, nil)
	_ = c(ctx, cmd)
	return cmd
}

// SPopN Redis `SPOP key count` command.
func (c cmdable) SPopN(ctx context.Context, key string, count int64) *StringSliceCmd {
	cmd := NewStringSliceCmd2(ctx, "spop", key, []interface{}{count})
	_ = c(ctx, cmd)
	return cmd
}

// SRandMember Redis `SRANDMEMBER key` command.
func (c cmdable) SRandMember(ctx context.Context, key string) *StringCmd {
	cmd := NewStringCmd2(ctx, "srandmember", key, nil)
	_ = c(ctx, cmd)
	return cmd
}

// SRandMemberN Redis `SRANDMEMBER key count` command.
func (c cmdable) SRandMemberN(ctx context.Context, key string, count int64) *StringSliceCmd {
	cmd := NewStringSliceCmd2(ctx, "srandmember", key, []interface{}{count})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) SRem(ctx context.Context, key string, members ...interface{}) *IntCmd {
	cmd := NewIntCmd2Any(ctx, "srem", key, members)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) SUnion(ctx context.Context, keys ...string) *StringSliceCmd {
	cmd := NewStringSliceCmd2S(ctx, "sunion", "", keys)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) SUnionStore(ctx context.Context, destination string, keys ...string) *IntCmd {
	cmd := NewIntCmd2S(ctx, "sunionstore", destination, keys)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd {
	args := []interface{}{cursor}
	if match != "" {
		args = append(args, "match", match)
	}
	if count > 0 {
		args = append(args, "count", count)
	}
	cmd := NewScanCmd2(ctx, c, "sscan", key, args)
	_ = c(ctx, cmd)
	return cmd
}
