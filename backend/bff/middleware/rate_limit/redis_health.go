package MiddlewareRateLimit

import "sync/atomic"

type RedisHealth struct {
	alive atomic.Bool
}

func (r *RedisHealth) SetAlive(v bool) {
	r.alive.Store(v)
}

func (r *RedisHealth) IsAlive() bool {
	return r.alive.Load()
}

func NewRedisHealth() *RedisHealth {
	return &RedisHealth{}
}
