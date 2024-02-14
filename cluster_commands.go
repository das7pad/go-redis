package redis

import "context"

type ClusterCmdable interface {
	ClusterMyShardID(ctx context.Context) *StringCmd
	ClusterSlots(ctx context.Context) *ClusterSlotsCmd
	ClusterShards(ctx context.Context) *ClusterShardsCmd
	ClusterLinks(ctx context.Context) *ClusterLinksCmd
	ClusterNodes(ctx context.Context) *StringCmd
	ClusterMeet(ctx context.Context, host, port string) *StatusCmd
	ClusterForget(ctx context.Context, nodeID string) *StatusCmd
	ClusterReplicate(ctx context.Context, nodeID string) *StatusCmd
	ClusterResetSoft(ctx context.Context) *StatusCmd
	ClusterResetHard(ctx context.Context) *StatusCmd
	ClusterInfo(ctx context.Context) *StringCmd
	ClusterKeySlot(ctx context.Context, key string) *IntCmd
	ClusterGetKeysInSlot(ctx context.Context, slot int, count int) *StringSliceCmd
	ClusterCountFailureReports(ctx context.Context, nodeID string) *IntCmd
	ClusterCountKeysInSlot(ctx context.Context, slot int) *IntCmd
	ClusterDelSlots(ctx context.Context, slots ...int) *StatusCmd
	ClusterDelSlotsRange(ctx context.Context, min, max int) *StatusCmd
	ClusterSaveConfig(ctx context.Context) *StatusCmd
	ClusterSlaves(ctx context.Context, nodeID string) *StringSliceCmd
	ClusterFailover(ctx context.Context) *StatusCmd
	ClusterAddSlots(ctx context.Context, slots ...int) *StatusCmd
	ClusterAddSlotsRange(ctx context.Context, min, max int) *StatusCmd
	ReadOnly(ctx context.Context) *StatusCmd
	ReadWrite(ctx context.Context) *StatusCmd
}

func (c cmdable) ClusterMyShardID(ctx context.Context) *StringCmd {
	cmd := NewStringCmd2(ctx, "cluster", "myshardid", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterSlots(ctx context.Context) *ClusterSlotsCmd {
	cmd := NewClusterSlotsCmd2(ctx, "cluster", "slots", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterShards(ctx context.Context) *ClusterShardsCmd {
	cmd := NewClusterShardsCmd2(ctx, "cluster", "shards", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterLinks(ctx context.Context) *ClusterLinksCmd {
	cmd := NewClusterLinksCmd2(ctx, "cluster", "links", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterNodes(ctx context.Context) *StringCmd {
	cmd := NewStringCmd2(ctx, "cluster", "nodes", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterMeet(ctx context.Context, host, port string) *StatusCmd {
	cmd := NewStatusCmd3S(ctx, "cluster", "meet", host, []string{port})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterForget(ctx context.Context, nodeID string) *StatusCmd {
	cmd := NewStatusCmd3(ctx, "cluster", "forget", nodeID, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterReplicate(ctx context.Context, nodeID string) *StatusCmd {
	cmd := NewStatusCmd3(ctx, "cluster", "replicate", nodeID, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterResetSoft(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd3(ctx, "cluster", "reset", "soft", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterResetHard(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd3(ctx, "cluster", "reset", "hard", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterInfo(ctx context.Context) *StringCmd {
	cmd := NewStringCmd2(ctx, "cluster", "info", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterKeySlot(ctx context.Context, key string) *IntCmd {
	cmd := NewIntCmd3(ctx, "cluster", "keyslot", key, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterGetKeysInSlot(ctx context.Context, slot int, count int) *StringSliceCmd {
	cmd := NewStringSliceCmd2(ctx, "cluster", "getkeysinslot", []interface{}{slot, count})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterCountFailureReports(ctx context.Context, nodeID string) *IntCmd {
	cmd := NewIntCmd3(ctx, "cluster", "count-failure-reports", nodeID, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterCountKeysInSlot(ctx context.Context, slot int) *IntCmd {
	cmd := NewIntCmd2(ctx, "cluster", "countkeysinslot", []interface{}{slot})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterDelSlots(ctx context.Context, slots ...int) *StatusCmd {
	args := make([]interface{}, len(slots))
	for i, slot := range slots {
		args[i] = slot
	}
	cmd := NewStatusCmd2(ctx, "cluster", "delslots", args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterDelSlotsRange(ctx context.Context, min, max int) *StatusCmd {
	size := max - min + 1
	slots := make([]int, size)
	for i := 0; i < size; i++ {
		slots[i] = min + i
	}
	return c.ClusterDelSlots(ctx, slots...)
}

func (c cmdable) ClusterSaveConfig(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "cluster", "saveconfig", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterSlaves(ctx context.Context, nodeID string) *StringSliceCmd {
	cmd := NewStringSliceCmd3(ctx, "cluster", "slaves", nodeID, nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterFailover(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "cluster", "failover", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterAddSlots(ctx context.Context, slots ...int) *StatusCmd {
	args := make([]interface{}, len(slots))
	for i, num := range slots {
		args[i] = num
	}
	cmd := NewStatusCmd2(ctx, "cluster", "addslots", args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ClusterAddSlotsRange(ctx context.Context, min, max int) *StatusCmd {
	size := max - min + 1
	slots := make([]int, size)
	for i := 0; i < size; i++ {
		slots[i] = min + i
	}
	return c.ClusterAddSlots(ctx, slots...)
}

func (c cmdable) ReadOnly(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "readonly", "", nil)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) ReadWrite(ctx context.Context) *StatusCmd {
	cmd := NewStatusCmd2(ctx, "readwrite", "", nil)
	_ = c(ctx, cmd)
	return cmd
}
