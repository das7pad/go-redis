package redis

import (
	"context"
	"encoding"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/redis/go-redis/v9/internal"
)

// KeepTTL is a Redis KEEPTTL option to keep existing TTL, it requires your redis-server version >= 6.0,
// otherwise you will receive an error: (error) ERR syntax error.
// For example:
//
//	rdb.Set(ctx, key, value, redis.KeepTTL)
const KeepTTL = -1

func usePrecise(dur time.Duration) bool {
	return dur < time.Second || dur%time.Second != 0
}

func formatMs(ctx context.Context, dur time.Duration) int64 {
	if dur > 0 && dur < time.Millisecond {
		internal.Logger.Printf(
			ctx,
			"specified duration is %s, but minimal supported value is %s - truncating to 1ms",
			dur, time.Millisecond,
		)
		return 1
	}
	return int64(dur / time.Millisecond)
}

func formatSec(ctx context.Context, dur time.Duration) int64 {
	if dur > 0 && dur < time.Second {
		internal.Logger.Printf(
			ctx,
			"specified duration is %s, but minimal supported value is %s - truncating to 1s",
			dur, time.Second,
		)
		return 1
	}
	return int64(dur / time.Second)
}

func appendArgs(dst, src []interface{}) ([]interface{}, []string) {
	if len(src) == 1 {
		return appendArg(dst, src[0])
	}
	if dst == nil {
		return src, nil
	}

	dst = append(dst, src...)
	return dst, nil
}

func appendArg(dst []interface{}, arg interface{}) ([]interface{}, []string) {
	switch arg := arg.(type) {
	case []string:
		if len(dst) > 0 {
			for _, s := range arg {
				dst = append(dst, s)
			}
			return dst, nil
		} else {
			return nil, arg[:len(arg):len(arg)]
		}
	case []interface{}:
		if len(dst) > 0 {
			return append(dst, arg...), nil
		} else {
			return arg[:len(arg):len(arg)], nil
		}
	case map[string]interface{}:
		for k, v := range arg {
			dst = append(dst, k, v)
		}
		return dst, nil
	case map[string]string:
		if len(dst) > 0 {
			for k, v := range arg {
				dst = append(dst, k, v)
			}
			return dst, nil
		}
		argsS := make([]string, 0, len(arg)*2)
		for k, v := range arg {
			argsS = append(argsS, k, v)
		}
		return nil, argsS
	case int64, time.Time, time.Duration, encoding.BinaryMarshaler, net.IP:
		return append(dst, arg), nil
	default:
		// scan struct field
		v := reflect.ValueOf(arg)
		if v.Type().Kind() == reflect.Ptr {
			if v.IsNil() {
				// error: arg is not a valid object
				return dst, nil
			}
			v = v.Elem()
		}

		if v.Type().Kind() == reflect.Struct {
			return appendStructField(dst, v), nil
		}

		return append(dst, arg), nil
	}
}

// appendStructField appends the field and value held by the structure v to dst, and returns the appended dst.
func appendStructField(dst []interface{}, v reflect.Value) []interface{} {
	typ := v.Type()
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("redis")
		if tag == "" || tag == "-" {
			continue
		}
		name, opt, _ := strings.Cut(tag, ",")
		if name == "" {
			continue
		}

		field := v.Field(i)

		// miss field
		if omitEmpty(opt) && isEmptyValue(field) {
			continue
		}

		if field.CanInterface() {
			dst = append(dst, name, field.Interface())
		}
	}

	return dst
}

func omitEmpty(opt string) bool {
	for opt != "" {
		var name string
		name, opt, _ = strings.Cut(opt, ",")
		if name == "omitempty" {
			return true
		}
	}
	return false
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Pointer:
		return v.IsNil()
	}
	return false
}

type Cmdable interface {
	Pipeline() Pipeliner
	Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error)

	TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error)
	TxPipeline() Pipeliner

	Command(ctx context.Context) *CommandsInfoCmd
	CommandList(ctx context.Context, filter *FilterBy) *StringSliceCmd
	CommandGetKeys(ctx context.Context, commands ...interface{}) *StringSliceCmd
	CommandGetKeysAndFlags(ctx context.Context, commands ...interface{}) *KeyFlagsCmd
	ClientGetName(ctx context.Context) *StringCmd
	Echo(ctx context.Context, message interface{}) *StringCmd
	Ping(ctx context.Context) *StatusCmd
	Quit(ctx context.Context) *StatusCmd
	Unlink(ctx context.Context, keys ...string) *IntCmd

	BgRewriteAOF(ctx context.Context) *StatusCmd
	BgSave(ctx context.Context) *StatusCmd
	ClientKill(ctx context.Context, ipPort string) *StatusCmd
	ClientKillByFilter(ctx context.Context, keys ...string) *IntCmd
	ClientList(ctx context.Context) *StringCmd
	ClientInfo(ctx context.Context) *ClientInfoCmd
	ClientPause(ctx context.Context, dur time.Duration) *BoolCmd
	ClientUnpause(ctx context.Context) *BoolCmd
	ClientID(ctx context.Context) *IntCmd
	ClientUnblock(ctx context.Context, id int64) *IntCmd
	ClientUnblockWithError(ctx context.Context, id int64) *IntCmd
	ConfigGet(ctx context.Context, parameter string) *MapStringStringCmd
	ConfigResetStat(ctx context.Context) *StatusCmd
	ConfigSet(ctx context.Context, parameter, value string) *StatusCmd
	ConfigRewrite(ctx context.Context) *StatusCmd
	DBSize(ctx context.Context) *IntCmd
	FlushAll(ctx context.Context) *StatusCmd
	FlushAllAsync(ctx context.Context) *StatusCmd
	FlushDB(ctx context.Context) *StatusCmd
	FlushDBAsync(ctx context.Context) *StatusCmd
	Info(ctx context.Context, section ...string) *StringCmd
	LastSave(ctx context.Context) *IntCmd
	Save(ctx context.Context) *StatusCmd
	Shutdown(ctx context.Context) *StatusCmd
	ShutdownSave(ctx context.Context) *StatusCmd
	ShutdownNoSave(ctx context.Context) *StatusCmd
	SlaveOf(ctx context.Context, host, port string) *StatusCmd
	SlowLogGet(ctx context.Context, num int64) *SlowLogCmd
	Time(ctx context.Context) *TimeCmd
	DebugObject(ctx context.Context, key string) *StringCmd
	MemoryUsage(ctx context.Context, key string, samples ...int) *IntCmd

	ModuleLoadex(ctx context.Context, conf *ModuleLoadexConfig) *StringCmd

	ACLCmdable
	BitMapCmdable
	ClusterCmdable
	GearsCmdable
	GenericCmdable
	GeoCmdable
	HashCmdable
	HyperLogLogCmdable
	ListCmdable
	ProbabilisticCmdable
	PubSubCmdable
	ScriptingFunctionsCmdable
	SetCmdable
	SortedSetCmdable
	StringCmdable
	StreamCmdable
	TimeseriesCmdable
	JSONCmdable
}

type StatefulCmdable interface {
	Cmdable
	Auth(ctx context.Context, password string) *StatusCmd
	AuthACL(ctx context.Context, username, password string) *StatusCmd
	Select(ctx context.Context, index int) *StatusCmd
	SwapDB(ctx context.Context, index1, index2 int) *StatusCmd
	ClientSetName(ctx context.Context, name string) *BoolCmd
	ClientSetInfo(ctx context.Context, info LibraryInfo) *StatusCmd
	Hello(ctx context.Context, ver int, username, password, clientName string) *MapStringInterfaceCmd
}

var (
	_ Cmdable = (*Client)(nil)
	_ Cmdable = (*Tx)(nil)
	_ Cmdable = (*Ring)(nil)
	_ Cmdable = (*ClusterClient)(nil)
)

type cmdable func(ctx context.Context, cmd Cmder) error

type statefulCmdable func(ctx context.Context, cmd Cmder) error

// ------------------------------------------------------------------------------

func (c statefulCmdable) Auth(ctx context.Context, password string) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "auth", password, nil)
	_ = c(ctx, cmd)
	return cmd
}

// AuthACL Perform an AUTH command, using the given user and pass.
// Should be used to authenticate the current connection with one of the connections defined in the ACL list
// when connecting to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
func (c statefulCmdable) AuthACL(ctx context.Context, username, password string) *StatusCmd {
	cmd := NewStatusCmd3(ctx, "auth", username, password, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) Wait(ctx context.Context, numSlaves int, timeout time.Duration) *IntCmd {
	cmd := NewIntCmd2(ctx, "wait", "", []interface{}{numSlaves, int(timeout / time.Millisecond)})
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) WaitAOF(ctx context.Context, numLocal, numSlaves int, timeout time.Duration) *IntCmd {
	cmd := NewIntCmd2(ctx, "waitAOF", "", []interface{}{numLocal, numSlaves, int(timeout / time.Millisecond)})
	cmd.setReadTimeout(timeout)
	_ = c(ctx, cmd)
	return cmd
}

func (c statefulCmdable) Select(ctx context.Context, index int) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "select", "", []interface{}{index})
	_ = c(ctx, cmd)
	return cmd
}

func (c statefulCmdable) SwapDB(ctx context.Context, index1, index2 int) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "swapdb", "", []interface{}{index1, index2})
	_ = c(ctx, cmd)
	return cmd
}

// ClientSetName assigns a name to the connection.
func (c statefulCmdable) ClientSetName(ctx context.Context, name string) *BoolCmd {
	cmd := NewBoolCmd3(ctx, "client", "setname", name, nil)
	_ = c(ctx, cmd)
	return cmd
}

// ClientSetInfo sends a CLIENT SETINFO command with the provided info.
func (c statefulCmdable) ClientSetInfo(ctx context.Context, info LibraryInfo) *StatusCmd {
	err := info.Validate()
	if err != nil {
		panic(err.Error())
	}

	var cmd *StatusCmd
	if info.LibName != nil {
		libName := fmt.Sprintf("go-redis(%s,%s)", *info.LibName, internal.ReplaceSpaces(runtime.Version()))
		cmd = NewStatusCmd3S(ctx, "client", "setinfo", "LIB-NAME", []string{libName})
	} else {
		cmd = NewStatusCmd3S(ctx, "client", "setinfo", "LIB-VER", []string{*info.LibVer})
	}

	_ = c(ctx, cmd)
	return cmd
}

// Validate checks if only one field in the struct is non-nil.
func (info LibraryInfo) Validate() error {
	if info.LibName != nil && info.LibVer != nil {
		return errors.New("both LibName and LibVer cannot be set at the same time")
	}
	if info.LibName == nil && info.LibVer == nil {
		return errors.New("at least one of LibName and LibVer should be set")
	}
	return nil
}

// Hello Set the resp protocol used.
func (c statefulCmdable) Hello(ctx context.Context,
	ver int, username, password, clientName string,
) *MapStringInterfaceCmd {
	args := make([]interface{}, 0, 6)
	args = append(args, ver)
	if password != "" {
		if username != "" {
			args = append(args, "auth", username, password)
		} else {
			args = append(args, "auth", "default", password)
		}
	}
	if clientName != "" {
		args = append(args, "setname", clientName)
	}
	cmd := NewMapStringInterfaceCmd2(ctx, "hello", "", args)
	_ = c(ctx, cmd)
	return cmd
}

// ------------------------------------------------------------------------------

func (c cmdable) Command(ctx context.Context) *CommandsInfoCmd {
	cmd := NewCommandsInfoCmd2(ctx, "command", "", nil)
	_ = c(ctx, cmd)
	return cmd
}

// FilterBy is used for the `CommandList` command parameter.
type FilterBy struct {
	Module  string
	ACLCat  string
	Pattern string
}

func (c cmdable) CommandList(ctx context.Context, filter *FilterBy) *StringSliceCmd {
	args := make([]string, 0, 3)
	if filter != nil {
		if filter.Module != "" {
			args = append(args, "filterby", "module", filter.Module)
		} else if filter.ACLCat != "" {
			args = append(args, "filterby", "aclcat", filter.ACLCat)
		} else if filter.Pattern != "" {
			args = append(args, "filterby", "pattern", filter.Pattern)
		}
	}
	cmd := NewStringSliceCmd2S(ctx, "command", "list", args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) CommandGetKeys(ctx context.Context, commands ...interface{}) *StringSliceCmd {
	cmd := NewStringSliceCmd2(ctx, "command", "getkeys", commands)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) CommandGetKeysAndFlags(ctx context.Context, commands ...interface{}) *KeyFlagsCmd {
	cmd := NewKeyFlagsCmd2(ctx, "command", "getkeysandflags", commands)
	_ = c(ctx, cmd)
	return cmd
}

// ClientGetName returns the name of the connection.
func (c cmdable) ClientGetName(ctx context.Context) *StringCmd {
	cmd := NewStringCmd2(ctx, "client", "getname", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) Echo(ctx context.Context, message interface{}) *StringCmd {
	cmd := NewStringCmd2(ctx, "echo", "", []interface{}{message})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) Ping(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "ping", "", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) Quit(_ context.Context) *StatusCmd {
	panic("not implemented")
}

// ------------------------------------------------------------------------------

func (c cmdable) BgRewriteAOF(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "bgrewriteaof", "", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) BgSave(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "bgsave", "", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClientKill(ctx context.Context, ipPort string) *StatusCmd {
	cmd := NewStatusCmd3(ctx, "client", "kill", ipPort, nil)
	_ = c(ctx, cmd)
	return cmd
}

// ClientKillByFilter is new style syntax, while the ClientKill is old
//
//	CLIENT KILL <option> [value] ... <option> [value]
func (c cmdable) ClientKillByFilter(ctx context.Context, keys ...string) *IntCmd {
	cmd := NewIntCmd2S(ctx, "client", "kill", keys)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClientList(ctx context.Context) *StringCmd {
	cmd := NewStringCmd2(ctx, "client", "list", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClientPause(ctx context.Context, dur time.Duration) *BoolCmd {
	cmd := NewBoolCmd2(ctx, "client", "pause", []interface{}{formatMs(ctx, dur)})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClientUnpause(ctx context.Context) *BoolCmd {
	cmd := NewBoolCmd2(ctx, "client", "unpause", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClientID(ctx context.Context) *IntCmd {
	cmd := NewIntCmd2(ctx, "client", "id", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClientUnblock(ctx context.Context, id int64) *IntCmd {
	cmd := NewIntCmd2(ctx, "client", "unblock", []interface{}{id})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClientUnblockWithError(ctx context.Context, id int64) *IntCmd {
	cmd := NewIntCmd2(ctx, "client", "unblock", []interface{}{id, "error"})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClientInfo(ctx context.Context) *ClientInfoCmd {
	cmd := NewClientInfoCmd2(ctx, "client", "info", nil)
	_ = c(ctx, cmd)
	return cmd
}

// ------------------------------------------------------------------------------------------------

func (c cmdable) ConfigGet(ctx context.Context, parameter string) *MapStringStringCmd {
	cmd := NewMapStringStringCmd3(ctx, "config", "get", parameter, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ConfigResetStat(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "config", "resetstat", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ConfigSet(ctx context.Context, parameter, value string) *StatusCmd {
	cmd := NewStatusCmd3S(ctx, "config", "set", parameter, []string{value})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ConfigRewrite(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "config", "rewrite", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) DBSize(ctx context.Context) *IntCmd {
	cmd := NewIntCmd2(ctx, "dbsize", "", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) FlushAll(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "flushall", "", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) FlushAllAsync(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "flushall", "async", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) FlushDB(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "flushdb", "", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) FlushDBAsync(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "flushdb", "async", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) Info(ctx context.Context, sections ...string) *StringCmd {
	cmd := NewStringCmd2S(ctx, "info", "", sections)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) InfoMap(ctx context.Context, sections ...string) *InfoCmd {
	cmd := NewInfoCmd2S(ctx, "info", "", sections)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) LastSave(ctx context.Context) *IntCmd {
	cmd := NewIntCmd2(ctx, "lastsave", "", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) Save(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "save", "", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) shutdown(ctx context.Context, modifier string) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "shutdown", modifier, nil)
	_ = c(ctx, cmd)
	if err := cmd.Err(); err != nil {
		if err == io.EOF {
			// Server quit as expected.
			cmd.err = nil
		}
	} else {
		// Server did not quit. String reply contains the reason.
		cmd.err = errors.New(cmd.val)
		cmd.val = ""
	}
	return cmd
}

func (c cmdable) Shutdown(ctx context.Context) *StatusCmd {
	return c.shutdown(ctx, "")
}

func (c cmdable) ShutdownSave(ctx context.Context) *StatusCmd {
	return c.shutdown(ctx, "save")
}

func (c cmdable) ShutdownNoSave(ctx context.Context) *StatusCmd {
	return c.shutdown(ctx, "nosave")
}

func (c cmdable) SlaveOf(ctx context.Context, host, port string) *StatusCmd {
	cmd := NewStatusCmd3(ctx, "slaveof", host, port, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) SlowLogGet(ctx context.Context, num int64) *SlowLogCmd {
	cmd := NewSlowLogCmd2(context.Background(), "slowlog", "get", []interface{}{num})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) Sync(_ context.Context) {
	panic("not implemented")
}

func (c cmdable) Time(ctx context.Context) *TimeCmd {
	cmd := NewTimeCmd2(ctx, "time", "", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) DebugObject(ctx context.Context, key string) *StringCmd {
	cmd := NewStringCmd3(ctx, "debug", "object", key, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) MemoryUsage(ctx context.Context, key string, samples ...int) *IntCmd {
	var args []interface{}
	if len(samples) > 0 {
		if len(samples) != 1 {
			panic("MemoryUsage expects single sample count")
		}
		args = append(args, "SAMPLES", samples[0])
	}
	cmd := NewIntCmd3(ctx, "memory", "usage", key, args)
	cmd.SetFirstKeyPos(2)
	_ = c(ctx, cmd)
	return cmd
}

// ------------------------------------------------------------------------------

// ModuleLoadexConfig struct is used to specify the arguments for the MODULE LOADEX command of redis.
// `MODULE LOADEX path [CONFIG name value [CONFIG name value ...]] [ARGS args [args ...]]`
type ModuleLoadexConfig struct {
	Path string
	Conf map[string]interface{}
	Args []interface{}
}

func (c *ModuleLoadexConfig) toArgs() []interface{} {
	args := make([]interface{}, 3, 3+len(c.Conf)*3+len(c.Args)*2)
	args[0] = "MODULE"
	args[1] = "LOADEX"
	args[2] = c.Path
	for k, v := range c.Conf {
		args = append(args, "CONFIG", k, v)
	}
	for _, arg := range c.Args {
		args = append(args, "ARGS", arg)
	}
	return args
}

// ModuleLoadex Redis `MODULE LOADEX path [CONFIG name value [CONFIG name value ...]] [ARGS args [args ...]]` command.
func (c cmdable) ModuleLoadex(ctx context.Context, conf *ModuleLoadexConfig) *StringCmd {
	cmd := NewStringCmd(ctx, conf.toArgs()...)
	_ = c(ctx, cmd)
	return cmd
}

/*
Monitor - represents a Redis MONITOR command, allowing the user to capture
and process all commands sent to a Redis server. This mimics the behavior of
MONITOR in the redis-cli.

Notes:
- Using MONITOR blocks the connection to the server for itself. It needs a dedicated connection
- The user should create a channel of type string
- This runs concurrently in the background. Trigger via the Start and Stop functions
See further: Redis MONITOR command: https://redis.io/commands/monitor
*/
func (c cmdable) Monitor(ctx context.Context, ch chan string) *MonitorCmd {
	cmd := newMonitorCmd(ctx, ch)
	_ = c(ctx, cmd)
	return cmd
}
