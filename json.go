package redis

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/redis/go-redis/v9/internal/proto"
	"github.com/redis/go-redis/v9/internal/util"
)

// -------------------------------------------

type JSONCmdable interface {
	JSONArrAppend(ctx context.Context, key, path string, values ...interface{}) *IntSliceCmd
	JSONArrIndex(ctx context.Context, key, path string, value ...interface{}) *IntSliceCmd
	JSONArrIndexWithArgs(ctx context.Context, key, path string, options *JSONArrIndexArgs, value ...interface{}) *IntSliceCmd
	JSONArrInsert(ctx context.Context, key, path string, index int64, values ...interface{}) *IntSliceCmd
	JSONArrLen(ctx context.Context, key, path string) *IntSliceCmd
	JSONArrPop(ctx context.Context, key, path string, index int) *StringSliceCmd
	JSONArrTrim(ctx context.Context, key, path string) *IntSliceCmd
	JSONArrTrimWithArgs(ctx context.Context, key, path string, options *JSONArrTrimArgs) *IntSliceCmd
	JSONClear(ctx context.Context, key, path string) *IntCmd
	JSONDebugMemory(ctx context.Context, key, path string) *IntCmd
	JSONDel(ctx context.Context, key, path string) *IntCmd
	JSONForget(ctx context.Context, key, path string) *IntCmd
	JSONGet(ctx context.Context, key string, paths ...string) *JSONCmd
	JSONGetWithArgs(ctx context.Context, key string, options *JSONGetArgs, paths ...string) *JSONCmd
	JSONMerge(ctx context.Context, key, path string, value string) *StatusCmd
	JSONMSetArgs(ctx context.Context, docs []JSONSetArgs) *StatusCmd
	JSONMSet(ctx context.Context, params ...interface{}) *StatusCmd
	JSONMGet(ctx context.Context, path string, keys ...string) *JSONSliceCmd
	JSONNumIncrBy(ctx context.Context, key, path string, value float64) *JSONCmd
	JSONObjKeys(ctx context.Context, key, path string) *SliceCmd
	JSONObjLen(ctx context.Context, key, path string) *IntPointerSliceCmd
	JSONSet(ctx context.Context, key, path string, value interface{}) *StatusCmd
	JSONSetMode(ctx context.Context, key, path string, value interface{}, mode string) *StatusCmd
	JSONStrAppend(ctx context.Context, key, path, value string) *IntPointerSliceCmd
	JSONStrLen(ctx context.Context, key, path string) *IntPointerSliceCmd
	JSONToggle(ctx context.Context, key, path string) *IntPointerSliceCmd
	JSONType(ctx context.Context, key, path string) *JSONSliceCmd
}

type JSONSetArgs struct {
	Key   string
	Path  string
	Value interface{}
}

type JSONArrIndexArgs struct {
	Start int
	Stop  *int
}

type JSONArrTrimArgs struct {
	Start int
	Stop  *int
}

type JSONCmd struct {
	baseCmd
	val      string
	expanded interface{}
}

var _ Cmder = (*JSONCmd)(nil)

func newJSONCmd2Both(ctx context.Context, cmd, firstArg string, args []interface{}, argsS []string) *JSONCmd {
	return &JSONCmd{
		baseCmd: baseCmd{
			ctx:      ctx,
			cmd:      cmd,
			firstArg: firstArg,
			args:     args,
			argsS:    argsS,
		},
	}
}

func (cmd *JSONCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *JSONCmd) SetVal(val string) {
	cmd.val = val
}

func (cmd *JSONCmd) Val() string {
	if len(cmd.val) == 0 && cmd.expanded != nil {
		val, err := json.Marshal(cmd.expanded)
		if err != nil {
			cmd.SetErr(err)
			return ""
		}
		return string(val)

	} else {
		return cmd.val
	}
}

func (cmd *JSONCmd) Result() (string, error) {
	return cmd.Val(), cmd.Err()
}

func (cmd *JSONCmd) Expanded() (interface{}, error) {
	if len(cmd.val) != 0 && cmd.expanded == nil {
		err := json.Unmarshal([]byte(cmd.val), &cmd.expanded)
		if err != nil {
			return nil, err
		}
	}

	return cmd.expanded, nil
}

func (cmd *JSONCmd) readReply(rd *proto.Reader) error {
	// nil response from JSON.(M)GET (cmd.baseCmd.err will be "redis: nil")
	if cmd.baseCmd.Err() == Nil {
		cmd.val = ""
		return Nil
	}

	if readType, err := rd.PeekReplyType(); err != nil {
		return err
	} else if readType == proto.RespArray {

		size, err := rd.ReadArrayLen()
		if err != nil {
			return err
		}

		expanded := make([]interface{}, size)

		for i := 0; i < size; i++ {
			if expanded[i], err = rd.ReadReply(); err != nil {
				return err
			}
		}
		cmd.expanded = expanded

	} else {
		if str, err := rd.ReadString(); err != nil && err != Nil {
			return err
		} else if str == "" || err == Nil {
			cmd.val = ""
		} else {
			cmd.val = str
		}
	}

	return nil
}

// -------------------------------------------

type JSONSliceCmd struct {
	baseCmd
	val []interface{}
}

func NewJSONSliceCmd(ctx context.Context, args ...interface{}) *JSONSliceCmd {
	return &JSONSliceCmd{
		baseCmd: baseCmd{
			ctx:  ctx,
			args: args,
		},
	}
}

func NewJSONSliceCmd2S(ctx context.Context, cmd, firstArg string, argsS []string) *JSONSliceCmd {
	return &JSONSliceCmd{
		baseCmd: baseCmd{
			ctx:      ctx,
			cmd:      cmd,
			firstArg: firstArg,
			argsS:    argsS,
		},
	}
}

func NewJSONSliceCmd3(ctx context.Context, cmd, firstArg, secondArg string, args []interface{}) *JSONSliceCmd {
	return &JSONSliceCmd{
		baseCmd: baseCmd{
			ctx:       ctx,
			cmd:       cmd,
			firstArg:  firstArg,
			secondArg: secondArg,
			args:      args,
		},
	}
}

func (cmd *JSONSliceCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *JSONSliceCmd) SetVal(val []interface{}) {
	cmd.val = val
}

func (cmd *JSONSliceCmd) Val() []interface{} {
	return cmd.val
}

func (cmd *JSONSliceCmd) Result() ([]interface{}, error) {
	return cmd.val, cmd.err
}

func (cmd *JSONSliceCmd) readReply(rd *proto.Reader) error {
	if cmd.baseCmd.Err() == Nil {
		cmd.val = nil
		return Nil
	}

	if readType, err := rd.PeekReplyType(); err != nil {
		return err
	} else if readType == proto.RespArray {
		response, err := rd.ReadReply()
		if err != nil {
			return nil
		} else {
			cmd.val = response.([]interface{})
		}

	} else {
		n, err := rd.ReadArrayLen()
		if err != nil {
			return err
		}
		cmd.val = make([]interface{}, n)
		for i := 0; i < len(cmd.val); i++ {
			switch s, err := rd.ReadString(); {
			case err == Nil:
				cmd.val[i] = ""
			case err != nil:
				return err
			default:
				cmd.val[i] = s
			}
		}
	}
	return nil
}

/*******************************************************************************
*
* IntPointerSliceCmd
* used to represent a RedisJSON response where the result is either an integer or nil
*
*******************************************************************************/

type IntPointerSliceCmd struct {
	baseCmd
	val []*int64
}

// NewIntPointerSliceCmd initialises an IntPointerSliceCmd
func NewIntPointerSliceCmd(ctx context.Context, args ...interface{}) *IntPointerSliceCmd {
	return &IntPointerSliceCmd{
		baseCmd: baseCmd{
			ctx:  ctx,
			args: args,
		},
	}
}

func NewIntPointerSliceCmd3(ctx context.Context, cmd, firstArg, secondArg string, args []interface{}) *IntPointerSliceCmd {
	return &IntPointerSliceCmd{
		baseCmd: baseCmd{
			ctx:       ctx,
			cmd:       cmd,
			firstArg:  firstArg,
			secondArg: secondArg,
			args:      args,
		},
	}
}

func (cmd *IntPointerSliceCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *IntPointerSliceCmd) SetVal(val []*int64) {
	cmd.val = val
}

func (cmd *IntPointerSliceCmd) Val() []*int64 {
	return cmd.val
}

func (cmd *IntPointerSliceCmd) Result() ([]*int64, error) {
	return cmd.val, cmd.err
}

func (cmd *IntPointerSliceCmd) readReply(rd *proto.Reader) error {
	n, err := rd.ReadArrayLen()
	if err != nil {
		return err
	}
	cmd.val = make([]*int64, n)

	for i := 0; i < len(cmd.val); i++ {
		val, err := rd.ReadInt()
		if err != nil && err != Nil {
			return err
		} else if err != Nil {
			cmd.val[i] = &val
		}
	}

	return nil
}

// ------------------------------------------------------------------------------

// JSONArrAppend adds the provided JSON values to the end of the array at the given path.
// For more information, see https://redis.io/commands/json.arrappend
func (c cmdable) JSONArrAppend(ctx context.Context, key, path string, values ...interface{}) *IntSliceCmd {
	cmd := NewIntSliceCmd3(ctx, "JSON.ARRAPPEND", key, path, values)
	_ = c(ctx, cmd)
	return cmd
}

// JSONArrIndex searches for the first occurrence of the provided JSON value in the array at the given path.
// For more information, see https://redis.io/commands/json.arrindex
func (c cmdable) JSONArrIndex(ctx context.Context, key, path string, value ...interface{}) *IntSliceCmd {
	cmd := NewIntSliceCmd3(ctx, "JSON.ARRINDEX", key, path, value)
	_ = c(ctx, cmd)
	return cmd
}

// JSONArrIndexWithArgs searches for the first occurrence of a JSON value in an array while allowing the start and
// stop options to be provided.
// For more information, see https://redis.io/commands/json.arrindex
func (c cmdable) JSONArrIndexWithArgs(ctx context.Context, key, path string, options *JSONArrIndexArgs, value ...interface{}) *IntSliceCmd {
	if options != nil {
		value = append(value, options.Start)
		if options.Stop != nil {
			value = append(value, *options.Stop)
		}
	}
	cmd := NewIntSliceCmd3(ctx, "JSON.ARRINDEX", key, path, value)
	_ = c(ctx, cmd)
	return cmd
}

// JSONArrInsert inserts the JSON values into the array at the specified path before the index (shifts to the right).
// For more information, see https://redis.io/commands/json.arrinsert
func (c cmdable) JSONArrInsert(ctx context.Context, key, path string, index int64, values ...interface{}) *IntSliceCmd {
	args := make([]interface{}, 0, 1+len(values))
	args = append(args, index)
	args = append(args, values...)
	cmd := NewIntSliceCmd3(ctx, "JSON.ARRINSERT", key, path, args)
	_ = c(ctx, cmd)
	return cmd
}

// JSONArrLen reports the length of the JSON array at the specified path in the given key.
// For more information, see https://redis.io/commands/json.arrlen
func (c cmdable) JSONArrLen(ctx context.Context, key, path string) *IntSliceCmd {
	cmd := NewIntSliceCmd3(ctx, "JSON.ARRLEN", key, path, nil)
	_ = c(ctx, cmd)
	return cmd
}

// JSONArrPop removes and returns an element from the specified index in the array.
// For more information, see https://redis.io/commands/json.arrpop
func (c cmdable) JSONArrPop(ctx context.Context, key, path string, index int) *StringSliceCmd {
	cmd := NewStringSliceCmd3(ctx, "JSON.ARRPOP", key, path, []interface{}{index})
	_ = c(ctx, cmd)
	return cmd
}

// JSONArrTrim trims an array to contain only the specified inclusive range of elements.
// For more information, see https://redis.io/commands/json.arrtrim
func (c cmdable) JSONArrTrim(ctx context.Context, key, path string) *IntSliceCmd {
	cmd := NewIntSliceCmd3(ctx, "JSON.ARRTRIM", key, path, nil)
	_ = c(ctx, cmd)
	return cmd
}

// JSONArrTrimWithArgs trims an array to contain only the specified inclusive range of elements.
// For more information, see https://redis.io/commands/json.arrtrim
func (c cmdable) JSONArrTrimWithArgs(ctx context.Context, key, path string, options *JSONArrTrimArgs) *IntSliceCmd {
	var args []interface{}
	if options != nil {
		args = append(args, options.Start)
		if options.Stop != nil {
			args = append(args, *options.Stop)
		}
	}
	cmd := NewIntSliceCmd3(ctx, "JSON.ARRTRIM", key, path, args)
	_ = c(ctx, cmd)
	return cmd
}

// JSONClear clears container values (arrays/objects) and sets numeric values to 0.
// For more information, see https://redis.io/commands/json.clear
func (c cmdable) JSONClear(ctx context.Context, key, path string) *IntCmd {
	cmd := NewIntCmd3(ctx, "JSON.CLEAR", key, path, nil)
	_ = c(ctx, cmd)
	return cmd
}

// JSONDebugMemory reports a value's memory usage in bytes (unimplemented)
// For more information, see https://redis.io/commands/json.debug-memory
func (c cmdable) JSONDebugMemory(ctx context.Context, key, path string) *IntCmd {
	panic("not implemented")
}

// JSONDel deletes a value.
// For more information, see https://redis.io/commands/json.del
func (c cmdable) JSONDel(ctx context.Context, key, path string) *IntCmd {
	cmd := NewIntCmd3(ctx, "JSON.DEL", key, path, nil)
	_ = c(ctx, cmd)
	return cmd
}

// JSONForget deletes a value.
// For more information, see https://redis.io/commands/json.forget
func (c cmdable) JSONForget(ctx context.Context, key, path string) *IntCmd {
	cmd := NewIntCmd3(ctx, "JSON.FORGET", key, path, nil)
	_ = c(ctx, cmd)
	return cmd
}

// JSONGet returns the value at path in JSON serialized form. JSON.GET returns an
// array of strings. This function parses out the wrapping array but leaves the
// internal strings unprocessed by default (see Val())
// For more information - https://redis.io/commands/json.get/
func (c cmdable) JSONGet(ctx context.Context, key string, paths ...string) *JSONCmd {
	cmd := newJSONCmd2Both(ctx, "JSON.GET", key, nil, paths)
	_ = c(ctx, cmd)
	return cmd
}

type JSONGetArgs struct {
	Indent  string
	Newline string
	Space   string
}

// JSONGetWithArgs - Retrieves the value of a key from a JSON document.
// This function also allows for specifying additional options such as:
// Indention, NewLine and Space
// For more information - https://redis.io/commands/json.get/
func (c cmdable) JSONGetWithArgs(ctx context.Context, key string, options *JSONGetArgs, paths ...string) *JSONCmd {
	var args []interface{}
	if options != nil {
		if options.Indent != "" {
			args = append(args, "INDENT", options.Indent)
		}
		if options.Newline != "" {
			args = append(args, "NEWLINE", options.Newline)
		}
		if options.Space != "" {
			args = append(args, "SPACE", options.Space)
		}
		for _, path := range paths {
			args = append(args, path)
		}
	}
	cmd := newJSONCmd2Both(ctx, "JSON.GET", key, args, nil)
	_ = c(ctx, cmd)
	return cmd
}

// JSONMerge merges a given JSON value into matching paths.
// For more information, see https://redis.io/commands/json.merge
func (c cmdable) JSONMerge(ctx context.Context, key, path string, value string) *StatusCmd {
	cmd := NewStatusCmd3S(ctx, "JSON.MERGE", key, path, []string{value})
	_ = c(ctx, cmd)
	return cmd
}

// JSONMGet returns the values at the specified path from multiple key arguments.
// Note - the arguments are reversed when compared with `JSON.MGET` as we want
// to follow the pattern of having the last argument be variable.
// For more information, see https://redis.io/commands/json.mget
func (c cmdable) JSONMGet(ctx context.Context, path string, keys ...string) *JSONSliceCmd {
	cmd := NewJSONSliceCmd2S(ctx, "JSON.MGET", "", append(keys, path))
	_ = c(ctx, cmd)
	return cmd
}

// JSONMSetArgs sets or updates one or more JSON values according to the specified key-path-value triplets.
// For more information, see https://redis.io/commands/json.mset
func (c cmdable) JSONMSetArgs(ctx context.Context, docs []JSONSetArgs) *StatusCmd {
	args := make([]interface{}, 0, 3*len(docs))
	for _, doc := range docs {
		args = append(args, doc.Key, doc.Path, doc.Value)
	}
	cmd := NewStatusCmd2(ctx, "JSON.MSET", "", args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) JSONMSet(ctx context.Context, params ...interface{}) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "JSON.MSET", "", params)
	_ = c(ctx, cmd)
	return cmd
}

// JSONNumIncrBy increments the number value stored at the specified path by the provided number.
// For more information, see https://redis.io/docs/latest/commands/json.numincrby/
func (c cmdable) JSONNumIncrBy(ctx context.Context, key, path string, value float64) *JSONCmd {
	cmd := newJSONCmd2Both(ctx, "JSON.NUMINCRBY", key, []interface{}{path, value}, nil)
	_ = c(ctx, cmd)
	return cmd
}

// JSONObjKeys returns the keys in the object that's referenced by the specified path.
// For more information, see https://redis.io/commands/json.objkeys
func (c cmdable) JSONObjKeys(ctx context.Context, key, path string) *SliceCmd {
	cmd := NewSliceCmd3(ctx, "JSON.OBJKEYS", key, path, nil)
	_ = c(ctx, cmd)
	return cmd
}

// JSONObjLen reports the number of keys in the JSON object at the specified path in the given key.
// For more information, see https://redis.io/commands/json.objlen
func (c cmdable) JSONObjLen(ctx context.Context, key, path string) *IntPointerSliceCmd {
	cmd := NewIntPointerSliceCmd3(ctx, "JSON.OBJLEN", key, path, nil)
	_ = c(ctx, cmd)
	return cmd
}

// JSONSet sets the JSON value at the given path in the given key. The value must be something that
// can be marshaled to JSON (using encoding/JSON) unless the argument is a string or a []byte when we assume that
// it can be passed directly as JSON.
// For more information, see https://redis.io/commands/json.set
func (c cmdable) JSONSet(ctx context.Context, key, path string, value interface{}) *StatusCmd {
	return c.JSONSetMode(ctx, key, path, value, "")
}

// JSONSetMode sets the JSON value at the given path in the given key and allows the mode to be set
// (the mode value must be "XX" or "NX"). The value must be something that can be marshaled to JSON (using encoding/JSON) unless
// the argument is a string or []byte when we assume that it can be passed directly as JSON.
// For more information, see https://redis.io/commands/json.set
func (c cmdable) JSONSetMode(ctx context.Context, key, path string, value interface{}, mode string) *StatusCmd {
	var bytes []byte
	var err error
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		bytes, err = json.Marshal(v)
	}
	args := []string{util.BytesToString(bytes)}
	if mode != "" {
		switch strings.ToUpper(mode) {
		case "XX", "NX":
			args = append(args, strings.ToUpper(mode))

		default:
			panic("redis: JSON.SET mode must be NX or XX")
		}
	}
	cmd := NewStatusCmd3S(ctx, "JSON.SET", key, path, args)
	if err != nil {
		cmd.SetErr(err)
	} else {
		_ = c(ctx, cmd)
	}
	return cmd
}

// JSONStrAppend appends the JSON-string values to the string at the specified path.
// For more information, see https://redis.io/commands/json.strappend
func (c cmdable) JSONStrAppend(ctx context.Context, key, path, value string) *IntPointerSliceCmd {
	cmd := NewIntPointerSliceCmd3(ctx, "JSON.STRAPPEND", key, path, []interface{}{value})
	_ = c(ctx, cmd)
	return cmd
}

// JSONStrLen reports the length of the JSON String at the specified path in the given key.
// For more information, see https://redis.io/commands/json.strlen
func (c cmdable) JSONStrLen(ctx context.Context, key, path string) *IntPointerSliceCmd {
	cmd := NewIntPointerSliceCmd3(ctx, "JSON.STRLEN", key, path, nil)
	_ = c(ctx, cmd)
	return cmd
}

// JSONToggle toggles a Boolean value stored at the specified path.
// For more information, see https://redis.io/commands/json.toggle
func (c cmdable) JSONToggle(ctx context.Context, key, path string) *IntPointerSliceCmd {
	cmd := NewIntPointerSliceCmd3(ctx, "JSON.TOGGLE", key, path, nil)
	_ = c(ctx, cmd)
	return cmd
}

// JSONType reports the type of JSON value at the specified path.
// For more information, see https://redis.io/commands/json.type
func (c cmdable) JSONType(ctx context.Context, key, path string) *JSONSliceCmd {
	cmd := NewJSONSliceCmd3(ctx, "JSON.TYPE", key, path, nil)
	_ = c(ctx, cmd)
	return cmd
}
