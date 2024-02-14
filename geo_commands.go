package redis

import (
	"context"
	"errors"
)

type GeoCmdable interface {
	GeoAdd(ctx context.Context, key string, geoLocation ...*GeoLocation) *IntCmd
	GeoPos(ctx context.Context, key string, members ...string) *GeoPosCmd
	GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *GeoRadiusQuery) *GeoLocationCmd
	GeoRadiusStore(ctx context.Context, key string, longitude, latitude float64, query *GeoRadiusQuery) *IntCmd
	GeoRadiusByMember(ctx context.Context, key, member string, query *GeoRadiusQuery) *GeoLocationCmd
	GeoRadiusByMemberStore(ctx context.Context, key, member string, query *GeoRadiusQuery) *IntCmd
	GeoSearch(ctx context.Context, key string, q *GeoSearchQuery) *StringSliceCmd
	GeoSearchLocation(ctx context.Context, key string, q *GeoSearchLocationQuery) *GeoSearchLocationCmd
	GeoSearchStore(ctx context.Context, key, store string, q *GeoSearchStoreQuery) *IntCmd
	GeoDist(ctx context.Context, key string, member1, member2, unit string) *FloatCmd
	GeoHash(ctx context.Context, key string, members ...string) *StringSliceCmd
}

func (c cmdable) GeoAdd(ctx context.Context, key string, geoLocation ...*GeoLocation) *IntCmd {
	args := make([]interface{}, 3*len(geoLocation))
	for i, eachLoc := range geoLocation {
		args[3*i] = eachLoc.Longitude
		args[3*i+1] = eachLoc.Latitude
		args[3*i+2] = eachLoc.Name
	}
	cmd := NewIntCmd2(ctx, "geoadd", key, args)
	_ = c(ctx, cmd)
	return cmd
}

// GeoRadius is a read-only GEORADIUS_RO command.
func (c cmdable) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *GeoRadiusQuery) *GeoLocationCmd {
	args := make([]interface{}, 0, 13)
	args = append(args, longitude, latitude)
	cmd := NewGeoLocationCmd2(ctx, query, "georadius_ro", key, args)
	if query.Store != "" || query.StoreDist != "" {
		cmd.SetErr(errors.New("GeoRadius does not support Store or StoreDist"))
		return cmd
	}
	_ = c(ctx, cmd)
	return cmd
}

// GeoRadiusStore is a writing GEORADIUS command.
func (c cmdable) GeoRadiusStore(ctx context.Context, key string, longitude, latitude float64, query *GeoRadiusQuery) *IntCmd {
	args := make([]interface{}, 0, 13)
	args = append(args, longitude, latitude)
	args = geoLocationArgs(query, args)
	cmd := NewIntCmd2(ctx, "georadius", key, args)
	if query.Store == "" && query.StoreDist == "" {
		cmd.SetErr(errors.New("GeoRadiusStore requires Store or StoreDist"))
		return cmd
	}
	_ = c(ctx, cmd)
	return cmd
}

// GeoRadiusByMember is a read-only GEORADIUSBYMEMBER_RO command.
func (c cmdable) GeoRadiusByMember(ctx context.Context, key, member string, query *GeoRadiusQuery) *GeoLocationCmd {
	args := make([]interface{}, 0, 11)
	cmd := NewGeoLocationCmd3(ctx, query, "georadiusbymember_ro", key, member, args)
	if query.Store != "" || query.StoreDist != "" {
		cmd.SetErr(errors.New("GeoRadiusByMember does not support Store or StoreDist"))
		return cmd
	}
	_ = c(ctx, cmd)
	return cmd
}

// GeoRadiusByMemberStore is a writing GEORADIUSBYMEMBER command.
func (c cmdable) GeoRadiusByMemberStore(ctx context.Context, key, member string, query *GeoRadiusQuery) *IntCmd {
	args := geoLocationArgs(query, make([]interface{}, 0, 11))
	cmd := NewIntCmd3(ctx, "georadiusbymember", key, member, args)
	if query.Store == "" && query.StoreDist == "" {
		cmd.SetErr(errors.New("GeoRadiusByMemberStore requires Store or StoreDist"))
		return cmd
	}
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) GeoSearch(ctx context.Context, key string, q *GeoSearchQuery) *StringSliceCmd {
	args := geoSearchArgs(q, make([]interface{}, 0, 11))
	cmd := NewStringSliceCmd2(ctx, "geosearch", key, args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) GeoSearchLocation(ctx context.Context, key string, q *GeoSearchLocationQuery) *GeoSearchLocationCmd {
	args := geoSearchLocationArgs(q, make([]interface{}, 0, 14))
	cmd := NewGeoSearchLocationCmd2(ctx, q, "geosearch", key, args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) GeoSearchStore(ctx context.Context, key, store string, q *GeoSearchStoreQuery) *IntCmd {
	args := geoSearchArgs(&q.GeoSearchQuery, make([]interface{}, 0, 12))
	if q.StoreDist {
		args = append(args, "storedist")
	}
	cmd := NewIntCmd3(ctx, "geosearchstore", store, key, args)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) GeoDist(ctx context.Context, key string, member1, member2, unit string) *FloatCmd {
	if unit == "" {
		unit = "km"
	}
	cmd := NewFloatCmd3S(ctx, "geodist", key, member1, []string{member2, unit})
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) GeoHash(ctx context.Context, key string, members ...string) *StringSliceCmd {
	cmd := NewStringSliceCmd2S(ctx, "geohash", key, members)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) GeoPos(ctx context.Context, key string, members ...string) *GeoPosCmd {
	cmd := NewGeoPosCmd2S(ctx, "geopos", key, members)
	_ = c(ctx, cmd)
	return cmd
}
