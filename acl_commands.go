package redis

import "context"

type ACLCmdable interface {
	ACLDryRun(ctx context.Context, username string, command ...interface{}) *StringCmd
	ACLLog(ctx context.Context, count int64) *ACLLogCmd
	ACLLogReset(ctx context.Context) *StatusCmd
}

func (c cmdable) ACLDryRun(ctx context.Context, username string, command ...interface{}) *StringCmd {
	cmd := NewStringCmd3(ctx, "acl", "dryrun", username, command)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ACLLog(ctx context.Context, count int64) *ACLLogCmd {
	var args []interface{}
	if count > 0 {
		args = append(args, count)
	}
	cmd := NewACLLogCmd2(ctx, "acl", "log", args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ACLLogReset(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd3(ctx, "acl", "log", "reset", nil)
	_ = c(ctx, cmd)
	return cmd
}
