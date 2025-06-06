package client

import "time"

type Options struct {
	DialTimeOut time.Duration
}

type OptionFunc func(*Options)

func WithDialTimeOut(timeout int64) OptionFunc {
	return func(o *Options) {
		o.DialTimeOut = time.Duration(timeout)
	}
}
