package redis

import (
	"context"
	"strings"
	"time"
)

type SortedSetCmdable interface {
	BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) *ZWithKeyCmd
	BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) *ZWithKeyCmd
	BZMPop(ctx context.Context, timeout time.Duration, order string, count int64, keys ...string) *ZSliceWithKeyCmd
	ZAdd(ctx context.Context, key string, members ...Z) *IntCmd
	ZAddLT(ctx context.Context, key string, members ...Z) *IntCmd
	ZAddGT(ctx context.Context, key string, members ...Z) *IntCmd
	ZAddNX(ctx context.Context, key string, members ...Z) *IntCmd
	ZAddXX(ctx context.Context, key string, members ...Z) *IntCmd
	ZAddArgs(ctx context.Context, key string, args ZAddArgs) *IntCmd
	ZAddArgsIncr(ctx context.Context, key string, args ZAddArgs) *FloatCmd
	ZCard(ctx context.Context, key string) *IntCmd
	ZCount(ctx context.Context, key, min, max string) *IntCmd
	ZLexCount(ctx context.Context, key, min, max string) *IntCmd
	ZIncrBy(ctx context.Context, key string, increment float64, member string) *FloatCmd
	ZInter(ctx context.Context, store *ZStore) *StringSliceCmd
	ZInterWithScores(ctx context.Context, store *ZStore) *ZSliceCmd
	ZInterCard(ctx context.Context, limit int64, keys ...string) *IntCmd
	ZInterStore(ctx context.Context, destination string, store *ZStore) *IntCmd
	ZMPop(ctx context.Context, order string, count int64, keys ...string) *ZSliceWithKeyCmd
	ZMScore(ctx context.Context, key string, members ...string) *FloatSliceCmd
	ZPopMax(ctx context.Context, key string, count ...int64) *ZSliceCmd
	ZPopMin(ctx context.Context, key string, count ...int64) *ZSliceCmd
	ZRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd
	ZRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd
	ZRangeByLex(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd
	ZRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd
	ZRangeArgs(ctx context.Context, z ZRangeArgs) *StringSliceCmd
	ZRangeArgsWithScores(ctx context.Context, z ZRangeArgs) *ZSliceCmd
	ZRangeStore(ctx context.Context, dst string, z ZRangeArgs) *IntCmd
	ZRank(ctx context.Context, key, member string) *IntCmd
	ZRankWithScore(ctx context.Context, key, member string) *RankWithScoreCmd
	ZRem(ctx context.Context, key string, members ...interface{}) *IntCmd
	ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *IntCmd
	ZRemRangeByScore(ctx context.Context, key, min, max string) *IntCmd
	ZRemRangeByLex(ctx context.Context, key, min, max string) *IntCmd
	ZRevRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd
	ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd
	ZRevRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd
	ZRevRangeByLex(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd
	ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd
	ZRevRank(ctx context.Context, key, member string) *IntCmd
	ZRevRankWithScore(ctx context.Context, key, member string) *RankWithScoreCmd
	ZScore(ctx context.Context, key, member string) *FloatCmd
	ZUnionStore(ctx context.Context, dest string, store *ZStore) *IntCmd
	ZRandMember(ctx context.Context, key string, count int) *StringSliceCmd
	ZRandMemberWithScores(ctx context.Context, key string, count int) *ZSliceCmd
	ZUnion(ctx context.Context, store ZStore) *StringSliceCmd
	ZUnionWithScores(ctx context.Context, store ZStore) *ZSliceCmd
	ZDiff(ctx context.Context, keys ...string) *StringSliceCmd
	ZDiffWithScores(ctx context.Context, keys ...string) *ZSliceCmd
	ZDiffStore(ctx context.Context, destination string, keys ...string) *IntCmd
	ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd
}

// BZPopMax Redis `BZPOPMAX key [key ...] timeout` command.
func (c cmdable) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) *ZWithKeyCmd {
	args := make([]interface{}, len(keys)+1)
	for i, key := range keys {
		args[i] = key
	}
	args[len(args)-1] = formatSec(ctx, timeout)
	cmd := NewZWithKeyCmd2(ctx, "bzpopmax", "", args)
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
}

// BZPopMin Redis `BZPOPMIN key [key ...] timeout` command.
func (c cmdable) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) *ZWithKeyCmd {
	args := make([]interface{}, len(keys)+1)
	for i, key := range keys {
		args[i] = key
	}
	args[len(args)-1] = formatSec(ctx, timeout)
	cmd := NewZWithKeyCmd2(ctx, "bzpopmin", "", args)
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
}

// BZMPop is the blocking variant of ZMPOP.
// When any of the sorted sets contains elements, this command behaves exactly like ZMPOP.
// When all sorted sets are empty, Redis will block the connection until another client adds members to one of the keys or until the timeout elapses.
// A timeout of zero can be used to block indefinitely.
// example: client.BZMPop(ctx, 0,"max", 1, "set")
func (c cmdable) BZMPop(ctx context.Context, timeout time.Duration, order string, count int64, keys ...string) *ZSliceWithKeyCmd {
	args := make([]interface{}, 2+len(keys), 5+len(keys))
	args[0] = formatSec(ctx, timeout)
	args[1] = len(keys)
	for i, key := range keys {
		args[2+i] = key
	}
	args = append(args, strings.ToLower(order), "count", count)
	cmd := NewZSliceWithKeyCmd2(ctx, "bzmpop", "", args)
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
}

// ZAddArgs WARN: The GT, LT and NX options are mutually exclusive.
type ZAddArgs struct {
	NX      bool
	XX      bool
	LT      bool
	GT      bool
	Ch      bool
	Members []Z
}

func (c cmdable) zAddArgs(args ZAddArgs, incr bool) []interface{} {
	a := make([]interface{}, 0, 4+2*len(args.Members))

	// The GT, LT and NX options are mutually exclusive.
	if args.NX {
		a = append(a, "nx")
	} else {
		if args.XX {
			a = append(a, "xx")
		}
		if args.GT {
			a = append(a, "gt")
		} else if args.LT {
			a = append(a, "lt")
		}
	}
	if args.Ch {
		a = append(a, "ch")
	}
	if incr {
		a = append(a, "incr")
	}
	for _, m := range args.Members {
		a = append(a, m.Score)
		a = append(a, m.Member)
	}
	return a
}

func (c cmdable) ZAddArgs(ctx context.Context, key string, args ZAddArgs) *IntCmd {
	cmd := NewIntCmd2(ctx, "zadd", key, c.zAddArgs(args, false))
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZAddArgsIncr(ctx context.Context, key string, args ZAddArgs) *FloatCmd {
	cmd := NewFloatCmd2(ctx, "zadd", key, c.zAddArgs(args, true))
	_ = c(ctx, cmd)
	return cmd
}

// ZAdd Redis `ZADD key score member [score member ...]` command.
func (c cmdable) ZAdd(ctx context.Context, key string, members ...Z) *IntCmd {
	return c.ZAddArgs(ctx, key, ZAddArgs{
		Members: members,
	})
}

// ZAddLT Redis `ZADD key LT score member [score member ...]` command.
func (c cmdable) ZAddLT(ctx context.Context, key string, members ...Z) *IntCmd {
	return c.ZAddArgs(ctx, key, ZAddArgs{
		LT:      true,
		Members: members,
	})
}

// ZAddGT Redis `ZADD key GT score member [score member ...]` command.
func (c cmdable) ZAddGT(ctx context.Context, key string, members ...Z) *IntCmd {
	return c.ZAddArgs(ctx, key, ZAddArgs{
		GT:      true,
		Members: members,
	})
}

// ZAddNX Redis `ZADD key NX score member [score member ...]` command.
func (c cmdable) ZAddNX(ctx context.Context, key string, members ...Z) *IntCmd {
	return c.ZAddArgs(ctx, key, ZAddArgs{
		NX:      true,
		Members: members,
	})
}

// ZAddXX Redis `ZADD key XX score member [score member ...]` command.
func (c cmdable) ZAddXX(ctx context.Context, key string, members ...Z) *IntCmd {
	return c.ZAddArgs(ctx, key, ZAddArgs{
		XX:      true,
		Members: members,
	})
}

func (c cmdable) ZCard(ctx context.Context, key string) *IntCmd {
	cmd := NewIntCmd2(ctx, "zcard", key, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZCount(ctx context.Context, key, min, max string) *IntCmd {
	cmd := NewIntCmd3S(ctx, "zcount", key, min, []string{max})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZLexCount(ctx context.Context, key, min, max string) *IntCmd {
	cmd := NewIntCmd3S(ctx, "zlexcount", key, min, []string{max})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZIncrBy(ctx context.Context, key string, increment float64, member string) *FloatCmd {
	cmd := NewFloatCmd2(ctx, "zincrby", key, []interface{}{increment, member})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZInterStore(ctx context.Context, destination string, store *ZStore) *IntCmd {
	args := store.appendArgs(make([]interface{}, 0, store.len()))
	cmd := NewIntCmd2(ctx, "zinterstore", destination, args)
	cmd.SetFirstKeyPos(3)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZInter(ctx context.Context, store *ZStore) *StringSliceCmd {
	args := store.appendArgs(make([]interface{}, 0, store.len()))
	cmd := NewStringSliceCmd2(ctx, "zinter", "", args)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZInterWithScores(ctx context.Context, store *ZStore) *ZSliceCmd {
	args := store.appendArgs(make([]interface{}, 0, store.len()+1))
	args = append(args, "withscores")
	cmd := NewZSliceCmd2(ctx, "zinter", "", args)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZInterCard(ctx context.Context, limit int64, keys ...string) *IntCmd {
	args := make([]interface{}, 1+len(keys), 1+len(keys)+2)
	args[0] = len(keys)
	for i, key := range keys {
		args[i+1] = key
	}
	args = append(args, "limit", limit)
	cmd := NewIntCmd2(ctx, "zintercard", "", args)
	_ = c(ctx, cmd)
	return cmd
}

// ZMPop Pops one or more elements with the highest or lowest score from the first non-empty sorted set key from the list of provided key names.
// direction: "max" (highest score) or "min" (lowest score), count: > 0
// example: client.ZMPop(ctx, "max", 5, "set1", "set2")
func (c cmdable) ZMPop(ctx context.Context, order string, count int64, keys ...string) *ZSliceWithKeyCmd {
	args := make([]interface{}, 1+len(keys), 1+len(keys)+3)
	args[0] = len(keys)
	for i, key := range keys {
		args[1+i] = key
	}
	args = append(args, strings.ToLower(order), "count", count)
	cmd := NewZSliceWithKeyCmd2(ctx, "zmpop", "", args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZMScore(ctx context.Context, key string, members ...string) *FloatSliceCmd {
	cmd := NewFloatSliceCmd2S(ctx, "zmscore", key, members)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZPopMax(ctx context.Context, key string, count ...int64) *ZSliceCmd {
	var args []interface{}
	switch len(count) {
	case 0:
		break
	case 1:
		args = append(args, count[0])
	default:
		panic("too many arguments")
	}

	cmd := NewZSliceCmd2(ctx, "zpopmax", key, args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZPopMin(ctx context.Context, key string, count ...int64) *ZSliceCmd {
	var args []interface{}
	switch len(count) {
	case 0:
		break
	case 1:
		args = append(args, count[0])
	default:
		panic("too many arguments")
	}

	cmd := NewZSliceCmd2(ctx, "zpopmin", key, args)
	_ = c(ctx, cmd)
	return cmd
}

// ZRangeArgs is all the options of the ZRange command.
// In version> 6.2.0, you can replace the(cmd):
//
//	ZREVRANGE,
//	ZRANGEBYSCORE,
//	ZREVRANGEBYSCORE,
//	ZRANGEBYLEX,
//	ZREVRANGEBYLEX.
//
// Please pay attention to your redis-server version.
//
// Rev, ByScore, ByLex and Offset+Count options require redis-server 6.2.0 and higher.
type ZRangeArgs struct {
	Key string

	// When the ByScore option is provided, the open interval(exclusive) can be set.
	// By default, the score intervals specified by <Start> and <Stop> are closed (inclusive).
	// It is similar to the deprecated(6.2.0+) ZRangeByScore command.
	// For example:
	//		ZRangeArgs{
	//			Key: 				"example-key",
	//	 		Start: 				"(3",
	//	 		Stop: 				8,
	//			ByScore:			true,
	//	 	}
	// 	 	cmd: "ZRange example-key (3 8 ByScore"  (3 < score <= 8).
	//
	// For the ByLex option, it is similar to the deprecated(6.2.0+) ZRangeByLex command.
	// You can set the <Start> and <Stop> options as follows:
	//		ZRangeArgs{
	//			Key: 				"example-key",
	//	 		Start: 				"[abc",
	//	 		Stop: 				"(def",
	//			ByLex:				true,
	//	 	}
	//		cmd: "ZRange example-key [abc (def ByLex"
	//
	// For normal cases (ByScore==false && ByLex==false), <Start> and <Stop> should be set to the index range (int).
	// You can read the documentation for more information: https://redis.io/commands/zrange
	Start interface{}
	Stop  interface{}

	// The ByScore and ByLex options are mutually exclusive.
	ByScore bool
	ByLex   bool

	Rev bool

	// limit offset count.
	Offset int64
	Count  int64
}

func (z ZRangeArgs) appendArgs(args []interface{}) []interface{} {
	// For Rev+ByScore/ByLex, we need to adjust the position of <Start> and <Stop>.
	if z.Rev && (z.ByScore || z.ByLex) {
		args = append(args, z.Key, z.Stop, z.Start)
	} else {
		args = append(args, z.Key, z.Start, z.Stop)
	}

	if z.ByScore {
		args = append(args, "byscore")
	} else if z.ByLex {
		args = append(args, "bylex")
	}
	if z.Rev {
		args = append(args, "rev")
	}
	if z.Offset != 0 || z.Count != 0 {
		args = append(args, "limit", z.Offset, z.Count)
	}
	return args
}

func (c cmdable) ZRangeArgs(ctx context.Context, z ZRangeArgs) *StringSliceCmd {
	args := z.appendArgs(make([]interface{}, 0, 8))
	cmd := NewStringSliceCmd2(ctx, "zrange", "", args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRangeArgsWithScores(ctx context.Context, z ZRangeArgs) *ZSliceCmd {
	args := z.appendArgs(make([]interface{}, 0, 9))
	args = append(args, "withscores")
	cmd := NewZSliceCmd2(ctx, "zrange", "", args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd {
	return c.ZRangeArgs(ctx, ZRangeArgs{
		Key:   key,
		Start: start,
		Stop:  stop,
	})
}

func (c cmdable) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd {
	return c.ZRangeArgsWithScores(ctx, ZRangeArgs{
		Key:   key,
		Start: start,
		Stop:  stop,
	})
}

type ZRangeBy struct {
	Min, Max      string
	Offset, Count int64
}

func (c cmdable) zRangeBy(ctx context.Context, zcmd, key string, opt *ZRangeBy, withScores bool) *StringSliceCmd {
	args := make([]interface{}, 0, 6)
	args = append(args, opt.Min, opt.Max)
	if withScores {
		args = append(args, "withscores")
	}
	if opt.Offset != 0 || opt.Count != 0 {
		args = append(args, "limit", opt.Offset, opt.Count)
	}
	cmd := NewStringSliceCmd2(ctx, zcmd, key, args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd {
	return c.zRangeBy(ctx, "zrangebyscore", key, opt, false)
}

func (c cmdable) ZRangeByLex(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd {
	return c.zRangeBy(ctx, "zrangebylex", key, opt, false)
}

func (c cmdable) ZRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd {
	args := make([]interface{}, 0, 6)
	args = append(args, opt.Min, opt.Max, "withscores")
	if opt.Offset != 0 || opt.Count != 0 {
		args = append(args, "limit", opt.Offset, opt.Count)
	}
	cmd := NewZSliceCmd2(ctx, "zrangebyscore", key, args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRangeStore(ctx context.Context, dst string, z ZRangeArgs) *IntCmd {
	args := z.appendArgs(make([]interface{}, 0, 8))
	cmd := NewIntCmd2(ctx, "zrangestore", dst, args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRank(ctx context.Context, key, member string) *IntCmd {
	cmd := NewIntCmd3(ctx, "zrank", key, member, nil)
	_ = c(ctx, cmd)
	return cmd
}

// ZRankWithScore according to the Redis documentation, if member does not exist
// in the sorted set or key does not exist, it will return a redis.Nil error.
func (c cmdable) ZRankWithScore(ctx context.Context, key, member string) *RankWithScoreCmd {
	cmd := NewRankWithScoreCmd3S(ctx, "zrank", key, member, []string{"withscore"})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRem(ctx context.Context, key string, members ...interface{}) *IntCmd {
	cmd := NewIntCmd2Any(ctx, "zrem", key, members)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *IntCmd {
	cmd := NewIntCmd2(ctx, "zremrangebyrank", key, []interface{}{start, stop})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRemRangeByScore(ctx context.Context, key, min, max string) *IntCmd {
	cmd := NewIntCmd2(ctx, "zremrangebyscore", key, []interface{}{min, max})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRemRangeByLex(ctx context.Context, key, min, max string) *IntCmd {
	cmd := NewIntCmd2(ctx, "zremrangebylex", key, []interface{}{min, max})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRevRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd {
	cmd := NewStringSliceCmd2(ctx, "zrevrange", key, []interface{}{start, stop})
	_ = c(ctx, cmd)
	return cmd
}

// ZRevRangeWithScores according to the Redis documentation, if member does not exist
// in the sorted set or key does not exist, it will return a redis.Nil error.
func (c cmdable) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *ZSliceCmd {
	cmd := NewZSliceCmd2(ctx, "zrevrange", key, []interface{}{start, stop, "withscores"})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) zRevRangeBy(ctx context.Context, zcmd, key string, opt *ZRangeBy) *StringSliceCmd {
	args := make([]interface{}, 0, 5)
	args = append(args, opt.Max, opt.Min)
	if opt.Offset != 0 || opt.Count != 0 {
		args = append(args, "limit", opt.Offset, opt.Count)
	}
	cmd := NewStringSliceCmd2(ctx, zcmd, key, args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRevRangeByScore(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd {
	return c.zRevRangeBy(ctx, "zrevrangebyscore", key, opt)
}

func (c cmdable) ZRevRangeByLex(ctx context.Context, key string, opt *ZRangeBy) *StringSliceCmd {
	return c.zRevRangeBy(ctx, "zrevrangebylex", key, opt)
}

func (c cmdable) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) *ZSliceCmd {
	args := []interface{}{opt.Max, opt.Min, "withscores"}
	if opt.Offset != 0 || opt.Count != 0 {
		args = append(args, "limit", opt.Offset, opt.Count)
	}
	cmd := NewZSliceCmd2(ctx, "zrevrangebyscore", key, args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRevRank(ctx context.Context, key, member string) *IntCmd {
	cmd := NewIntCmd3(ctx, "zrevrank", key, member, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZRevRankWithScore(ctx context.Context, key, member string) *RankWithScoreCmd {
	cmd := NewRankWithScoreCmd3S(ctx, "zrevrank", key, member, []string{"withscore"})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZScore(ctx context.Context, key, member string) *FloatCmd {
	cmd := NewFloatCmd3(ctx, "zscore", key, member, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZUnion(ctx context.Context, store ZStore) *StringSliceCmd {
	args := store.appendArgs(make([]interface{}, 0, store.len()))
	cmd := NewStringSliceCmd2(ctx, "zunion", "", args)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZUnionWithScores(ctx context.Context, store ZStore) *ZSliceCmd {
	args := store.appendArgs(make([]interface{}, 0, store.len()+1))
	args = append(args, "withscores")
	cmd := NewZSliceCmd2(ctx, "zunion", "", args)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZUnionStore(ctx context.Context, dest string, store *ZStore) *IntCmd {
	args := store.appendArgs(make([]interface{}, 0, store.len()))
	cmd := NewIntCmd2(ctx, "zunionstore", dest, args)
	cmd.SetFirstKeyPos(3)
	_ = c(ctx, cmd)
	return cmd
}

// ZRandMember redis-server version >= 6.2.0.
func (c cmdable) ZRandMember(ctx context.Context, key string, count int) *StringSliceCmd {
	cmd := NewStringSliceCmd2(ctx, "zrandmember", key, []interface{}{count})
	_ = c(ctx, cmd)
	return cmd
}

// ZRandMemberWithScores redis-server version >= 6.2.0.
func (c cmdable) ZRandMemberWithScores(ctx context.Context, key string, count int) *ZSliceCmd {
	cmd := NewZSliceCmd2(ctx, "zrandmember", key, []interface{}{count, "withscores"})
	_ = c(ctx, cmd)
	return cmd
}

// ZDiff redis-server version >= 6.2.0.
func (c cmdable) ZDiff(ctx context.Context, keys ...string) *StringSliceCmd {
	args := make([]interface{}, 1+len(keys))
	args[0] = len(keys)
	for i, key := range keys {
		args[i+1] = key
	}

	cmd := NewStringSliceCmd2(ctx, "zdiff", "", args)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
}

// ZDiffWithScores redis-server version >= 6.2.0.
func (c cmdable) ZDiffWithScores(ctx context.Context, keys ...string) *ZSliceCmd {
	args := make([]interface{}, 1+len(keys), 1+len(keys)+1)
	args[0] = len(keys)
	for i, key := range keys {
		args[i+1] = key
	}
	args = append(args, "withscores")

	cmd := NewZSliceCmd2(ctx, "zdiff", "", args)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
}

// ZDiffStore redis-server version >=6.2.0.
func (c cmdable) ZDiffStore(ctx context.Context, destination string, keys ...string) *IntCmd {
	args := make([]interface{}, 1+len(keys), 1+len(keys)+1)
	args[0] = len(keys)
	for i, key := range keys {
		args[i+1] = key
	}
	cmd := NewIntCmd2(ctx, "zdiffstore", destination, args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd {
	args := []interface{}{cursor}
	if match != "" {
		args = append(args, "match", match)
	}
	if count > 0 {
		args = append(args, "count", count)
	}
	cmd := NewScanCmd2(ctx, c, "zscan", key, args)
	_ = c(ctx, cmd)
	return cmd
}

// Z represents sorted set member.
type Z struct {
	Score  float64
	Member interface{}
}

// ZWithKey represents sorted set member including the name of the key where it was popped.
type ZWithKey struct {
	Z
	Key string
}

// ZStore is used as an arg to ZInter/ZInterStore and ZUnion/ZUnionStore.
type ZStore struct {
	Keys    []string
	Weights []float64
	// Can be SUM, MIN or MAX.
	Aggregate string
}

func (z ZStore) len() int {
	n := 1 // length of z.Keys
	n += len(z.Keys)
	if len(z.Weights) > 0 {
		n += 1 + len(z.Weights)
	}
	if z.Aggregate != "" {
		n += 2
	}
	return n
}

func (z ZStore) appendArgs(args []interface{}) []interface{} {
	args = append(args, len(z.Keys))
	for _, key := range z.Keys {
		args = append(args, key)
	}
	if len(z.Weights) > 0 {
		args = append(args, "weights")
		for _, weights := range z.Weights {
			args = append(args, weights)
		}
	}
	if z.Aggregate != "" {
		args = append(args, "aggregate", z.Aggregate)
	}
	return args
}
