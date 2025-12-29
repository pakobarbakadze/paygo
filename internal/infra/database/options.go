package database

import "time"

type Options struct {
	maxIdleConns    int
	maxOpenConns    int
	connMaxLifetime time.Duration
	timeout         time.Duration
	maxRetries      int
	retryDelay      time.Duration
}

func defaultOptions() *Options {
	return &Options{
		maxIdleConns:    10,
		maxOpenConns:    100,
		connMaxLifetime: time.Hour,
		timeout:         10 * time.Second,
		maxRetries:      3,
		retryDelay:      2 * time.Second,
	}
}

type Option func(*Options)

func WithConnectionPool(maxIdle, maxOpen int, maxLifetime time.Duration) Option {
	return func(o *Options) {
		o.maxIdleConns = maxIdle
		o.maxOpenConns = maxOpen
		o.connMaxLifetime = maxLifetime
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.timeout = timeout
	}
}

func WithMaxRetries(maxRetries int) Option {
	return func(o *Options) {
		o.maxRetries = maxRetries
	}
}

func WithRetryDelay(delay time.Duration) Option {
	return func(o *Options) {
		o.retryDelay = delay
	}
}
