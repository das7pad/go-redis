package redis

import (
	"context"
	"strings"
)

type GearsCmdable interface {
	TFunctionLoad(ctx context.Context, lib string) *StatusCmd
	TFunctionLoadArgs(ctx context.Context, lib string, options *TFunctionLoadOptions) *StatusCmd
	TFunctionDelete(ctx context.Context, libName string) *StatusCmd
	TFunctionList(ctx context.Context) *MapStringInterfaceSliceCmd
	TFunctionListArgs(ctx context.Context, options *TFunctionListOptions) *MapStringInterfaceSliceCmd
	TFCall(ctx context.Context, libName string, funcName string, numKeys int) *Cmd
	TFCallArgs(ctx context.Context, libName string, funcName string, numKeys int, options *TFCallOptions) *Cmd
	TFCallASYNC(ctx context.Context, libName string, funcName string, numKeys int) *Cmd
	TFCallASYNCArgs(ctx context.Context, libName string, funcName string, numKeys int, options *TFCallOptions) *Cmd
}

type TFunctionLoadOptions struct {
	Replace bool
	Config  string
}

type TFunctionListOptions struct {
	Withcode bool
	Verbose  int
	Library  string
}

type TFCallOptions struct {
	Keys      []string
	Arguments []string
}

// TFunctionLoad - load a new JavaScript library into Redis.
// For more information - https://redis.io/commands/tfunction-load/
func (c cmdable) TFunctionLoad(ctx context.Context, lib string) *StatusCmd {
	cmd := NewStatusCmd3(ctx, "TFUNCTION", "LOAD", lib, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) TFunctionLoadArgs(ctx context.Context, lib string, options *TFunctionLoadOptions) *StatusCmd {
	var argsS []string
	if options != nil {
		if options.Replace {
			argsS = append(argsS, "REPLACE")
		}
		if options.Config != "" {
			argsS = append(argsS, "CONFIG", options.Config)
		}
	}
	argsS = append(argsS, lib)
	cmd := NewStatusCmd2S(ctx, "TFUNCTION", "LOAD", argsS)
	_ = c(ctx, cmd)
	return cmd
}

// TFunctionDelete - delete a JavaScript library from Redis.
// For more information - https://redis.io/commands/tfunction-delete/
func (c cmdable) TFunctionDelete(ctx context.Context, libName string) *StatusCmd {
	cmd := NewStatusCmd3(ctx, "TFUNCTION", "DELETE", libName, nil)
	_ = c(ctx, cmd)
	return cmd
}

// TFunctionList - list the functions with additional information about each function.
// For more information - https://redis.io/commands/tfunction-list/
func (c cmdable) TFunctionList(ctx context.Context) *MapStringInterfaceSliceCmd {
	cmd := NewMapStringInterfaceSliceCmd2(ctx, "TFUNCTION", "LIST", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) TFunctionListArgs(ctx context.Context, options *TFunctionListOptions) *MapStringInterfaceSliceCmd {
	var args []interface{}
	if options != nil {
		if options.Withcode {
			args = append(args, "WITHCODE")
		}
		if options.Verbose != 0 {
			v := strings.Repeat("v", options.Verbose)
			args = append(args, v)
		}
		if options.Library != "" {
			args = append(args, "LIBRARY", options.Library)
		}
	}
	cmd := NewMapStringInterfaceSliceCmd2(ctx, "TFUNCTION", "LIST", args)
	_ = c(ctx, cmd)
	return cmd
}

// TFCall - invoke a function.
// For more information - https://redis.io/commands/tfcall/
func (c cmdable) TFCall(ctx context.Context, libName string, funcName string, numKeys int) *Cmd {
	cmd := NewCmd2(ctx, "TFCALL", libName+"."+funcName, []interface{}{numKeys})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) TFCallArgs(ctx context.Context, libName string, funcName string, numKeys int, options *TFCallOptions) *Cmd {
	args := []interface{}{numKeys}
	if options != nil {
		for _, key := range options.Keys {
			args = append(args, key)
		}
		for _, key := range options.Arguments {
			args = append(args, key)
		}
	}
	cmd := NewCmd2(ctx, "TFCALL", libName+"."+funcName, args)
	_ = c(ctx, cmd)
	return cmd
}

// TFCallASYNC - invoke an asynchronous JavaScript function (coroutine).
// For more information - https://redis.io/commands/TFCallASYNC/
func (c cmdable) TFCallASYNC(ctx context.Context, libName string, funcName string, numKeys int) *Cmd {
	cmd := NewCmd2(ctx, "TFCALLASYNC", libName+"."+funcName, []interface{}{numKeys})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) TFCallASYNCArgs(ctx context.Context, libName string, funcName string, numKeys int, options *TFCallOptions) *Cmd {
	args := []interface{}{numKeys}
	if options != nil {
		for _, key := range options.Keys {
			args = append(args, key)
		}
		for _, key := range options.Arguments {
			args = append(args, key)
		}
	}
	cmd := NewCmd2(ctx, "TFCALLASYNC", libName+"."+funcName, args)
	_ = c(ctx, cmd)
	return cmd
}
