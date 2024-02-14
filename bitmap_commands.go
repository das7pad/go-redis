package redis

import (
	"context"
	"errors"
)

type BitMapCmdable interface {
	GetBit(ctx context.Context, key string, offset int64) *IntCmd
	SetBit(ctx context.Context, key string, offset int64, value int) *IntCmd
	BitCount(ctx context.Context, key string, bitCount *BitCount) *IntCmd
	BitOpAnd(ctx context.Context, destKey string, keys ...string) *IntCmd
	BitOpOr(ctx context.Context, destKey string, keys ...string) *IntCmd
	BitOpXor(ctx context.Context, destKey string, keys ...string) *IntCmd
	BitOpNot(ctx context.Context, destKey string, key string) *IntCmd
	BitPos(ctx context.Context, key string, bit int64, pos ...int64) *IntCmd
	BitPosSpan(ctx context.Context, key string, bit int8, start, end int64, span string) *IntCmd
	BitField(ctx context.Context, key string, values ...interface{}) *IntSliceCmd
	BitFieldRO(ctx context.Context, key string, values ...interface{}) *IntSliceCmd
}

func (c cmdable) GetBit(ctx context.Context, key string, offset int64) *IntCmd {
	cmd := NewIntCmd2(ctx, "getbit", key, []interface{}{offset})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) SetBit(ctx context.Context, key string, offset int64, value int) *IntCmd {
	cmd := NewIntCmd2(ctx, "setbit", key, []interface{}{offset, value})
	_ = c(ctx, cmd)
	return cmd
}

type BitCount struct {
	Start, End int64
	Unit       string // BYTE(default) | BIT
}

const BitCountIndexByte string = "BYTE"
const BitCountIndexBit string = "BIT"

func (c cmdable) BitCount(ctx context.Context, key string, bitCount *BitCount) *IntCmd {
	var args []interface{}
	if bitCount != nil {
		args = append(args, bitCount.Start, bitCount.End)
		if bitCount.Unit != "" {
			if bitCount.Unit != BitCountIndexByte && bitCount.Unit != BitCountIndexBit {
				cmd := NewIntCmd(ctx)
				cmd.SetErr(errors.New("redis: invalid bitcount index"))
				return cmd
			}
			args = append(args, bitCount.Unit)
		}
	}
	cmd := NewIntCmd2(ctx, "bitcount", key, args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) bitOp(ctx context.Context, op, destKey string, keys []string) *IntCmd {
	cmd := NewIntCmd3S(ctx, "bitop", op, destKey, keys)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) BitOpAnd(ctx context.Context, destKey string, keys ...string) *IntCmd {
	return c.bitOp(ctx, "and", destKey, keys)
}

func (c cmdable) BitOpOr(ctx context.Context, destKey string, keys ...string) *IntCmd {
	return c.bitOp(ctx, "or", destKey, keys)
}

func (c cmdable) BitOpXor(ctx context.Context, destKey string, keys ...string) *IntCmd {
	return c.bitOp(ctx, "xor", destKey, keys)
}

func (c cmdable) BitOpNot(ctx context.Context, destKey string, key string) *IntCmd {
	return c.bitOp(ctx, "not", destKey, []string{key})
}

// BitPos is an API before Redis version 7.0, cmd: bitpos key bit start end
// if you need the `byte | bit` parameter, please use `BitPosSpan`.
func (c cmdable) BitPos(ctx context.Context, key string, bit int64, pos ...int64) *IntCmd {
	args := make([]interface{}, 0, 1+len(pos))
	args = append(args, bit)
	switch len(pos) {
	case 0:
	case 1:
		args = append(args, pos[0])
	case 2:
		args = append(args, pos[0])
		args = append(args, pos[1])
	default:
		panic("too many arguments")
	}
	cmd := NewIntCmd2(ctx, "bitpos", key, args)
	_ = c(ctx, cmd)
	return cmd
}

// BitPosSpan supports the `byte | bit` parameters in redis version 7.0,
// the bitpos command defaults to using byte type for the `start-end` range,
// which means it counts in bytes from start to end. you can set the value
// of "span" to determine the type of `start-end`.
// span = "bit", cmd: bitpos key bit start end bit
// span = "byte", cmd: bitpos key bit start end byte
func (c cmdable) BitPosSpan(ctx context.Context, key string, bit int8, start, end int64, span string) *IntCmd {
	cmd := NewIntCmd2(ctx, "bitpos", key, []interface{}{bit, start, end, span})
	_ = c(ctx, cmd)
	return cmd
}

// BitField accepts multiple values:
//   - BitField("set", "i1", "offset1", "value1","cmd2", "type2", "offset2", "value2")
//   - BitField([]string{"cmd1", "type1", "offset1", "value1","cmd2", "type2", "offset2", "value2"})
//   - BitField([]interface{}{"cmd1", "type1", "offset1", "value1","cmd2", "type2", "offset2", "value2"})
func (c cmdable) BitField(ctx context.Context, key string, values ...interface{}) *IntSliceCmd {
	cmd := NewIntSliceCmd2(ctx, "bitfield", key, values)
	_ = c(ctx, cmd)
	return cmd
}

// BitFieldRO - Read-only variant of the BITFIELD command.
// It is like the original BITFIELD but only accepts GET subcommand and can safely be used in read-only replicas.
// - BitFieldRO(ctx, key, "<Encoding0>", "<Offset0>", "<Encoding1>","<Offset1>")
func (c cmdable) BitFieldRO(ctx context.Context, key string, values ...interface{}) *IntSliceCmd {
	args := make([]interface{}, 0, len(values)+len(values)/2)
	if len(values)%2 != 0 {
		panic("BitFieldRO: invalid number of arguments, must be even")
	}
	for i := 0; i < len(values); i += 2 {
		args = append(args, "GET", values[i], values[i+1])
	}
	cmd := NewIntSliceCmd2(ctx, "BITFIELD_RO", key, args)
	_ = c(ctx, cmd)
	return cmd
}
