package server

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// PoolMetrics exposes Redis pool gauges using the global OTel meter provider.
type PoolMetrics struct {
	redisTotalConns metric.Int64Gauge
	redisIdleConns  metric.Int64Gauge
	redisHits       metric.Int64Gauge
	redisMisses     metric.Int64Gauge
}

func NewPoolMetrics(serviceName string) (*PoolMetrics, error) {
	meter := otel.Meter(serviceName)

	redisTotalConns, err := meter.Int64Gauge(
		"redis_pool_total_conns",
		metric.WithDescription("Total number of connections in the Redis pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	redisIdleConns, err := meter.Int64Gauge(
		"redis_pool_idle_conns",
		metric.WithDescription("Number of idle connections in the Redis pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	redisHits, err := meter.Int64Gauge(
		"redis_pool_hits_total",
		metric.WithDescription("Total number of times a free connection was found in the Redis pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	redisMisses, err := meter.Int64Gauge(
		"redis_pool_misses_total",
		metric.WithDescription("Total number of times a free connection was not found in the Redis pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	return &PoolMetrics{
		redisTotalConns: redisTotalConns,
		redisIdleConns:  redisIdleConns,
		redisHits:       redisHits,
		redisMisses:     redisMisses,
	}, nil
}

func (m *PoolMetrics) RecordRedisStats(ctx context.Context, total, idle, hits, misses int64) {
	m.redisTotalConns.Record(ctx, total)
	m.redisIdleConns.Record(ctx, idle)
	m.redisHits.Record(ctx, hits)
	m.redisMisses.Record(ctx, misses)
}
