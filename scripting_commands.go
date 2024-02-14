package redis

import "context"

type ScriptingFunctionsCmdable interface {
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *Cmd
	EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *Cmd
	EvalRO(ctx context.Context, script string, keys []string, args ...interface{}) *Cmd
	EvalShaRO(ctx context.Context, sha1 string, keys []string, args ...interface{}) *Cmd
	ScriptExists(ctx context.Context, hashes ...string) *BoolSliceCmd
	ScriptFlush(ctx context.Context) *StatusCmd
	ScriptKill(ctx context.Context) *StatusCmd
	ScriptLoad(ctx context.Context, script string) *StringCmd

	FunctionLoad(ctx context.Context, code string) *StringCmd
	FunctionLoadReplace(ctx context.Context, code string) *StringCmd
	FunctionDelete(ctx context.Context, libName string) *StringCmd
	FunctionFlush(ctx context.Context) *StringCmd
	FunctionKill(ctx context.Context) *StringCmd
	FunctionFlushAsync(ctx context.Context) *StringCmd
	FunctionList(ctx context.Context, q FunctionListQuery) *FunctionListCmd
	FunctionDump(ctx context.Context) *StringCmd
	FunctionRestore(ctx context.Context, libDump string) *StringCmd
	FunctionStats(ctx context.Context) *FunctionStatsCmd
	FCall(ctx context.Context, function string, keys []string, args ...interface{}) *Cmd
	FCallRo(ctx context.Context, function string, keys []string, args ...interface{}) *Cmd
	FCallRO(ctx context.Context, function string, keys []string, args ...interface{}) *Cmd
}

func (c cmdable) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *Cmd {
	return c.eval(ctx, "eval", script, keys, args...)
}

func (c cmdable) EvalRO(ctx context.Context, script string, keys []string, args ...interface{}) *Cmd {
	return c.eval(ctx, "eval_ro", script, keys, args...)
}

func (c cmdable) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *Cmd {
	return c.eval(ctx, "evalsha", sha1, keys, args...)
}

func (c cmdable) EvalShaRO(ctx context.Context, sha1 string, keys []string, args ...interface{}) *Cmd {
	return c.eval(ctx, "evalsha_ro", sha1, keys, args...)
}

func (c cmdable) eval(ctx context.Context, name, payload string, keys []string, args ...interface{}) *Cmd {
	cmd := NewCmd2(ctx, name, payload, fcallArgs(keys, args))

	// it is possible that only args exist without a key.
	// rdb.eval(ctx, eval, script, nil, arg1, arg2)
	if len(keys) > 0 {
		cmd.SetFirstKeyPos(3)
	}
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ScriptExists(ctx context.Context, hashes ...string) *BoolSliceCmd {
	cmd := NewBoolSliceCmd2S(ctx, "script", "exists", hashes)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ScriptFlush(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "script", "flush", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ScriptKill(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "script", "kill", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ScriptLoad(ctx context.Context, script string) *StringCmd {
	cmd := NewStringCmd3(ctx, "script", "load", script, nil)
	_ = c(ctx, cmd)
	return cmd
}

// ------------------------------------------------------------------------------

// FunctionListQuery is used with FunctionList to query for Redis libraries
//
//	  	LibraryNamePattern 	- Use an empty string to get all libraries.
//	  						- Use a glob-style pattern to match multiple libraries with a matching name
//	  						- Use a library's full name to match a single library
//		WithCode			- If true, it will return the code of the library
type FunctionListQuery struct {
	LibraryNamePattern string
	WithCode           bool
}

func (c cmdable) FunctionLoad(ctx context.Context, code string) *StringCmd {
	cmd := NewStringCmd3(ctx, "function", "load", code, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) FunctionLoadReplace(ctx context.Context, code string) *StringCmd {
	cmd := NewStringCmd3S(ctx, "function", "load", "replace", []string{code})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) FunctionDelete(ctx context.Context, libName string) *StringCmd {
	cmd := NewStringCmd3(ctx, "function", "delete", libName, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) FunctionFlush(ctx context.Context) *StringCmd {
	cmd := NewStringCmd2(ctx, "function", "flush", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) FunctionKill(ctx context.Context) *StringCmd {
	cmd := NewStringCmd2(ctx, "function", "kill", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) FunctionFlushAsync(ctx context.Context) *StringCmd {
	cmd := NewStringCmd3(ctx, "function", "flush", "async", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) FunctionList(ctx context.Context, q FunctionListQuery) *FunctionListCmd {
	argsS := make([]string, 0, 3)
	if q.LibraryNamePattern != "" {
		argsS = append(argsS, "libraryname", q.LibraryNamePattern)
	}
	if q.WithCode {
		argsS = append(argsS, "withcode")
	}
	cmd := NewFunctionListCmd2S(ctx, "function", "list", argsS)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) FunctionDump(ctx context.Context) *StringCmd {
	cmd := NewStringCmd2(ctx, "function", "dump", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) FunctionRestore(ctx context.Context, libDump string) *StringCmd {
	cmd := NewStringCmd3(ctx, "function", "restore", libDump, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) FunctionStats(ctx context.Context) *FunctionStatsCmd {
	cmd := NewFunctionStatsCmd2(ctx, "function", "stats", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) FCall(ctx context.Context, function string, keys []string, args ...interface{}) *Cmd {
	cmd := NewCmd2(ctx, "fcall", function, fcallArgs(keys, args))
	if len(keys) > 0 {
		cmd.SetFirstKeyPos(3)
	}
	_ = c(ctx, cmd)
	return cmd
}

// FCallRo this function simply calls FCallRO,
// Deprecated: to maintain convention FCallRO.
func (c cmdable) FCallRo(ctx context.Context, function string, keys []string, args ...interface{}) *Cmd {
	return c.FCallRO(ctx, function, keys, args...)
}

func (c cmdable) FCallRO(ctx context.Context, function string, keys []string, args ...interface{}) *Cmd {
	cmd := NewCmd2(ctx, "fcall_ro", function, fcallArgs(keys, args))
	if len(keys) > 0 {
		cmd.SetFirstKeyPos(3)
	}
	_ = c(ctx, cmd)
	return cmd
}

func fcallArgs(keys []string, args []interface{}) []interface{} {
	cmdArgs := make([]interface{}, 1+len(keys), 1+len(keys)+len(args))
	cmdArgs[0] = len(keys)
	for i, key := range keys {
		cmdArgs[1+i] = key
	}
	cmdArgs = append(cmdArgs, args...)
	return cmdArgs
}
